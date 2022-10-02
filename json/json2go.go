package main

import (
	"encoding/json"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/rs/zerolog/log"
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
	ApixRoots      string `json:"apixRoots"`
}

func GenerateJson2Go(content []byte) (string, error) {
	log.Info().Msgf("将json内容解析成go源码")
	data := model.ApixData{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Error().Msgf("无法将json转化为ApixData,请检查json内容是否符合格式，%v", err)
		return "", err
	} else {
		return GenerateGo(data)
	}
}

func Build(fileName string) (string, error) {
	goFile, err := GenerateFile2Go(fileName)
	if err != nil {
		return "", err
	}
	//编译
	buildMode, suffix := getBuildModeAndSuffix()
	//拼装pluginName
	pluginName := strings.ReplaceAll(goFile, ".go", "")
	log.Info().Msgf("开始构建插件：%s", pluginName)
	targetTmp := pluginName + "_tmp" + suffix
	target := pluginName + suffix
	if err = execBuild(buildMode, targetTmp, goFile); err != nil {
		return "", err
	}
	if err = execUpx(targetTmp, target); err != nil {
		//压缩失败，不影响使用
		target = targetTmp
	}
	return target, nil
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
		if e := os.Remove(filePath); e != nil {
			log.Warn().Msgf("goFile remove failed,%v", e)
		}

	}(goFile)
	buildCmd := exec.Command("go", "build", buildMode, "-o", targetTmp, goFile)
	err := buildCmd.Run()
	if err != nil {
		log.Error().Msgf("构建失败,cmd=%s:%v", buildCmd.String(), err.Error())
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

func GenerateGo(data model.ApixData) (string, error) {
	log.Info().Msgf("apixData对象解析成go源码")
	var roots []model.ApixStep
	for _, step := range data.Rule.Steps {
		if step.PrevId == "" {
			roots = append(roots, step)
		}
	}
	key := getKey(data.Rule.Api)
	tmpl, err := getTemplate(key)
	if err != nil {
		log.Error().Msgf("无法解析模板,%v", err)
		return "", err
	}
	tp := templateParam{
		ApixPath:       data.Rule.Api.Path,
		ApixParameters: obj2ByteArray(data.Rule.Api.Parameters),
		ApixResponse:   obj2ByteArray(data.Rule.Response),
		ApixRoots:      obj2ByteArray(roots),
		ApixSteps:      obj2ByteArray(data.Rule.Steps),
	}

	goFilePath := getGoFilePath(key)
	if f, err := createGoFile(goFilePath); err != nil {
		return "", err
	} else if err = tmpl.Execute(f, tp); err != nil {
		log.Error().Msgf("模板编译失败,%v,%s", err, debug.Stack())
		return "", err
	}
	return goFilePath, nil
}

func getGoFilePath(key string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, PLUGIN_PATH, key+".go")

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
		tmplateName := "json2go.tmpl"
		if runtime.GOOS == "windows" {
			tmplateName = "json2go_windows.tmpl"
		}
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

func getKey(api model.ApixApi) string {
	key := api.Code
	if key == "" || len(key) == 0 {
		key = strings.Join([]string{api.Path, api.Method, api.Version}, "_")
		key = strings.ReplaceAll(key, "/", "")
		if strings.HasPrefix(key, "_") {
			key = strings.ReplaceAll(key, "_", "bintools")
		}
	}
	if runtime.GOOS == "windows" {
		key += "_windows"
	}
	return key
}

func obj2ByteArray(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
