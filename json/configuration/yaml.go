package configuration

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

var config any

func init() {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "resources")
	defaultFile := filepath.Join(fp, "application.yml")
	if err := loadFile(defaultFile); err != nil {
		log.Fatal().Msgf("无法读取配置文件，:%v", err)
	}
	active := GetConfig("server.active")
	if active != nil {
		activeFile := filepath.Join(fp, fmt.Sprintf("application-%v.yml", active))
		if err := loadFile(activeFile); err != nil {
			log.Fatal().Msgf("无法读取配置文件，:%v", err)
		}
	}
}

func loadFile(fp string) error {
	if data, err := os.ReadFile(fp); err != nil {
		log.Error().Msgf("无法读取配置文件:[%s]:%v", fp, err)
	} else if err = yaml.Unmarshal(data, &config); err != nil {
		log.Error().Msgf("无法解析配置文件:[%s]:%v", fp, err)
	}
	return nil
}

func GetConfig(key string) any {
	return getConfigWithMap(key, config.(map[any]any))
}

func getConfigWithMap(key string, config map[any]any) any {
	if strings.Contains(key, ".") {
		keys := strings.Split(key, ".")
		subKey := keys[0]
		var subConfig map[any]any
		if v, ok := config[subKey]; ok {
			subConfig = v.(map[any]any)
			return getConfigWithMap(strings.Join(keys[1:], "."), subConfig)
		}
	}
	if v, ok := config[key]; ok {
		return v
	}
	return nil
}
