package json

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
	wd, _ := os.Getwd()
	shellPath := path.Join(wd, "plugins", "compile.sh")
	if runtime.GOOS == "windows" {
		shellPath = strings.ReplaceAll(shellPath, ".sh", ".bat")
	}
	pluginName := strings.ReplaceAll(goFile, ".go", "")
	output, err := exec.Command(shellPath, pluginName).Output()
	if err != nil {
		return "", err
	} else {
		log.Info().Msgf("构建成功:%s", output)
	}
	pluginName += ".so"
	if runtime.GOOS == "windows" {
		pluginName = strings.ReplaceAll(pluginName, ".so", ".dll")
	}
	return pluginName, nil
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
	//todo 给变量赋值{{steps.roots}}
	key := strings.Join([]string{data.Rule.Api.Path, data.Rule.Api.Method, data.Rule.Api.Version}, "_")
	key = strings.ReplaceAll(key, "/", "")
	t := template.New(key)
	templateData := func() string {
		wd, _ := os.Getwd()
		filePath := path.Join(wd, "template", "json2go.tmpl")
		if data, err := os.ReadFile(filePath); err != nil {
			log.Error().Msgf("无法读取模板内容，%v", err)
			return ""
		} else {
			return string(data)
		}
	}()
	tmpl, err := t.Parse(templateData)
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
	pwd, _ := os.Getwd()

	goFilePath := path.Join(pwd, "plugins", key+".go")
	f, err := os.Create(goFilePath)
	if err != nil {
		log.Error().Msgf("文件创建失败,%v", err)
		return "", err
	}
	err = tmpl.Execute(f, tp)
	if err != nil {
		log.Error().Msgf("模板编译失败,%v,%s", err, debug.Stack())
		return "", err
	}
	return key + ".go", nil
}

func obj2ByteArray(obj any) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
