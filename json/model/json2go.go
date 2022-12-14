package model

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
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
	TenantId       string `json:"tenantId"`
}

func GenerateJson2Go(content []byte, tenantId string) (string, error) {
	log.Info().Msgf("将json内容解析成apixData对象")
	data := ApixData{}
	if err := json.Unmarshal(content, &data); err != nil {
		log.Error().Msgf("无法将json转化为ApixData,请检查json内容是否符合格式，%v", err)
		return "", err
	} else {
		return generateGo(data, tenantId)
	}
}

//GenerateFile2Go 返回Go文件地址，或者错误信息
func GenerateFile2Go(fileName, tenantId string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Error().Msgf("无法打开文件,%v", err)
		return "", err
	}
	return GenerateJson2Go(data, tenantId)
}

func generateGo(data ApixData, tenantId string) (string, error) {
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
		TenantId:       tenantId,
	}

	goFilePath := getGoFilePath(pk, tenantId)
	if f, err := createGoFile(goFilePath); err != nil {
		return "", err
	} else if err = tmpl.Execute(f, tp); err != nil {
		log.Error().Msgf("模板编译失败,%v,%s", err, debug.Stack())
		return "", err
	}
	return goFilePath, nil
}

func getGoFilePath(key, tenantId string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "example", tenantId, key+".go_")

}

func createGoFile(goFilePath string) (*os.File, error) {
	idx := strings.LastIndex(goFilePath, string(os.PathSeparator))
	dirPath := goFilePath[0:idx]
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(dirPath, os.ModeDir)
		}
	}
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
