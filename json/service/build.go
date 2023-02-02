package service

import "github.com/kuchensheng/bintools/json/model"

func BuildJsonFile(filePath, tenantId, appCode string) (string, error) {
	goFile, err := model.GenerateFile2Go(filePath, tenantId, appCode)
	if err != nil {
		return "", err
	}
	err = Compile(goFile)
	return goFile, err //buildGoFile2Plugin(goFile)
}

func BuildJson(content []byte, tenantId, appCode string) (string, error) {
	if goFile, err := model.GenerateJson2Go(content, tenantId, appCode); err != nil {
		return "", err
	} else {
		go Compile(goFile)
		return goFile, err
	}
}
