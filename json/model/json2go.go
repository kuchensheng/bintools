package model

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"text/template"
)

const PLUGIN_PATH = "plugins"

type templateParam struct {
	ApixPath       string `json:"apixPath"`
	ApixParameters string `json:"apixParameters"`
	ApixResponse   string `json:"apixResponse"`
	ApixSteps      string `json:"apixSteps"`
	Key            string `json:"key"`
}

func GenerateJson2Go(content []byte) (string, error) {
	log.Info().Msgf("将json内容解析成apixData对象")
	data := ApixData{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Error().Msgf("无法将json转化为ApixData,请检查json内容是否符合格式，%v", err)
		return "", err
	} else {
		return GenerateGo(data)
	}
}

func BuildJsonFile(filePath string) (string, error) {
	goFile, err := GenerateFile2Go(filePath)
	if err != nil {
		return "", err
	}
	return goFile, err //buildGoFile2Plugin(goFile)
}

func BuildJson(content []byte) (string, error) {
	return GenerateJson2Go(content)
}

func Build(fileName string) (string, error) {
	return GenerateFile2Go(fileName)
}

func execUpx(targetTmp, target string) error {
	defer func() error {
		if x := recover(); x != nil {
			return x.(error)
		}
		return nil
	}()
	_ = os.Remove(target)
	upxCmd := exec.Command("upx", "-o", target, targetTmp)
	if err := upxCmd.Run(); err != nil {
		log.Error().Msgf("压缩失败%v", err)
		return err
	}
	log.Info().Msgf("压缩成功，plugin=%s", target)
	//删除临时文件
	removeTmpFile(targetTmp)
	return nil
}

func removeTmpFile(targetTmp string) {
	_ = os.Remove(targetTmp)
	hTargetTmp := strings.ReplaceAll(targetTmp, ".dll", ".h")
	_ = os.Remove(hTargetTmp)
}

func execBuild(buildMode, targetTmp, goFile string) error {
	defer func(filePath string) {
		if e := recover(); e != nil {
			log.Warn().Msgf("goFile build failed,%v", e)
		}
	}(goFile)
	if runtime.GOOS == "windows" {
		targetTmp = strings.ReplaceAll(targetTmp, "/", `\`)
		goFile = strings.ReplaceAll(goFile, "/", `\`)
	}
	buildCmd := exec.Command("go", "build", buildMode, "-o", targetTmp, goFile)
	cmdData, err := buildCmd.Output()
	if err != nil {
		log.Error().Msgf("构建失败,cmd=%s:%s", buildCmd.String(), debug.Stack())
		log.Error().Msgf("%s", cmdData)
		return err
	}

	log.Info().Msgf("构建成功:plugin=%s", targetTmp)
	return nil
}

func getBuildModeAndSuffix() (string, string) {
	buildMode := "-buildmode=plugin"
	suffix := ".so"
	if runtime.GOOS == "windows" {
		suffix = ".dll"
		buildMode = "-buildmode=c-shared"
	}
	return buildMode, suffix
}

func GenerateFile2Go(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Error().Msgf("无法打开文件,%v", err)
		return "", err
	}
	return GenerateJson2Go(data)
}

func GenerateGo(data ApixData) (string, error) {
	log.Info().Msgf("apixData对象解析成go源码")

	key := getKey(data.Rule.Api)
	tmpl, err := getTemplate(key)
	if err != nil {
		log.Error().Msgf("无法解析模板,%v", err)
		return "", err
	}
	pk := strings.ReplaceAll(key, "_api_app_orc_", "")
	pk = strings.ReplaceAll(pk, "_", "")
	tp := templateParam{
		ApixPath:       data.Rule.Api.Path,
		ApixParameters: obj2ByteArray(data.Rule.Api.Parameters),
		ApixResponse:   obj2ByteArray(data.Rule.Response),
		ApixSteps:      obj2ByteArray(data.Rule.Steps),
		Key:            pk,
	}

	goFilePath := getGoFilePath(pk)
	if f, err := createGoFile(goFilePath); err != nil {
		return "", err
	} else if err = tmpl.Execute(f, tp); err != nil {
		log.Error().Msgf("模板编译失败,%v,%s", err, debug.Stack())
		return "", err
	}
	go Compile(goFilePath, pk)
	return goFilePath, nil
}

func removeGoFile() {
	pwd, _ := os.Getwd()
	dir := path.Join(pwd, PLUGIN_PATH)
	readDir, _ := ioutil.ReadDir(dir)
	for _, info := range readDir {
		if strings.HasSuffix(info.Name(), "go") {
			if err := os.Remove(path.Join(dir, info.Name()+".go_")); err != nil {
				log.Error().Msgf("无法删除文件:%s", info.Name()+".go_,%v", err)
			}
		}
	}
}

func getGoFilePath(key string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "example", key+".go_")

}

func createGoFile(goFilePath string) (*os.File, error) {
	f, err := os.Create(goFilePath)
	if err != nil {
		log.Error().Msgf("文件创建失败,%v", err)
		return nil, err
	}
	return f, err
}

func getTemplate(key string) (*template.Template, error) {
	t := template.New(key)
	templateData := func() string {
		wd, _ := os.Getwd()
		tmplateName := "tmp.tmpl"
		filePath := path.Join(wd, "template", tmplateName)
		if data, err := os.ReadFile(filePath); err != nil {
			log.Error().Msgf("无法读取模板内容，%v", err)
			return ""
		} else {
			return string(data)
		}
	}()
	return t.Parse(templateData)
}

func getKey(api ApixApi) string {
	key := api.Code
	if key == "" || len(key) == 0 {
		method := api.Method
		if method == "" || len(method) == 0 {
			method = "get"
		} else {
			method = strings.ToLower(method)
		}
		key = strings.Join([]string{api.Path, method, api.Version}, "_")
		key = strings.ReplaceAll(key, "/", "_")
	}
	return key
}

func obj2ByteArray(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
