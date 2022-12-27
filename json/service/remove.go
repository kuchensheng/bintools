package service

import (
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/lib"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Remove(tenantId, path, method, json string) (int, error) {
	pk := strings.ReplaceAll(path, consts.GlobalPrefix, "")
	pk = strings.ReplaceAll(pk, "/", "")
	pk = pk + strings.ToLower(method)
	dir := filepath.Join(lib.Wd, "example", tenantId)
	goFile := filepath.Join(dir, pk+".go")
	if err := os.Remove(goFile); err != nil {
		log.Error().Msgf("无法删除go文件[%s],%v", goFile, err)
		return http.StatusInternalServerError, model.NewBusinessException(1080500, "无法删除Go脚本文件,文件不存在")
	}
	//移除执行引擎
	lib.RemoveProgramMap(pk)
	if json != "" {
		jsonFile := filepath.Join(dir, json)
		if err := os.Remove(jsonFile); err != nil {
			log.Error().Msgf("无法删除json文件[%s],%v", jsonFile, err)
			return http.StatusInternalServerError, model.NewBusinessException(1080500, "无法删除Json文件,"+err.Error())
		}
	}
	return http.StatusOK, model.NewBusinessException(0, "成功")
}
