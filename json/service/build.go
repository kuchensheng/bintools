package service

import "github.com/kuchensheng/bintools/json/model"

func BuildJsonFile(filePath, tenantId string) (string, error) {
	goFile, err := model.GenerateFile2Go(filePath, tenantId)
	if err != nil {
		return "", err
	}
	go Compile(goFile)
	return goFile, err //buildGoFile2Plugin(goFile)
}

func BuildJson(content []byte, tenantId string) (string, error) {
	return model.GenerateJson2Go(content, tenantId)
}
