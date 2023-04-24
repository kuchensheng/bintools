package yaml

import (
	"fmt"
	config2 "github.com/kuchensheng/bintools/http/config"
	"github.com/kuchensheng/bintools/logger"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var pwd = func() string {
	dir, _ := os.Getwd()
	return dir
}()
var (
	defaultYamlPath = filepath.Join(pwd, "resources", "application.yml")
)

type YamlConfig struct {
}

func InitYamlConfig() *YamlConfig {
	readYaml(defaultYamlPath)
	if active := getVal("spring.profiles.active", config2.ConfigMap); active != nil {
		activePath := filepath.Join(pwd, "resources", "application-"+active.(string)+".yml")
		readYaml(activePath)
	}
	return &YamlConfig{}
}

func readYaml(yamlPath string) {
	if yamlPath == "" {
		yamlPath = defaultYamlPath
	}
	file, err := os.ReadFile(yamlPath)
	if err != nil {
		logger.GlobalLogger.FatalLevel(fmt.Sprintf("无法读取配置文件,%v", err))
	}
	config := make(map[string]any)
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		logger.GlobalLogger.FatalLevel(fmt.Sprintf("无法解析配置文件,%v", err))
	}
	if config2.ConfigMap == nil {
		config2.ConfigMap = config
	} else {
		//同名覆盖Map
		updateConfig(config, config2.ConfigMap)
	}
}

func updateConfig(config map[string]any, oldConfigMap map[string]any) {
	for key, val := range config {
		if v, ok := oldConfigMap[key]; ok {
			if _, is := v.(map[string]any); is {
				updateConfig(val.(map[string]any), v.(map[string]any))
			} else {
				oldConfigMap[key] = v
			}
		} else {
			oldConfigMap[key] = val
		}
	}
}

func getVal(key string, m map[string]any) any {
	split := strings.Split(key, config2.Concat)
	if len(split) == 1 {
		return m[key]
	}
	subKey := split[0]
	if a, ok := m[subKey]; ok {
		if _, is := a.(map[string]any); !is {
			return a
		} else {
			return getVal(key[len(subKey)+1:], a.(map[string]any))
		}
	} else {
		return nil
	}
}

func (c *YamlConfig) GetAttr(key string) any {
	if config2.ConfigMap == nil {
		readYaml(defaultYamlPath)
	}
	return getVal(key, config2.ConfigMap)
}

func (c *YamlConfig) FillAttr(obj any, prefix string) {
	if obj == nil {
		return
	}
	if len(config2.ConfigMap) == 0 {
		readYaml(defaultYamlPath)
	}
	v1 := reflect.ValueOf(obj)
	t1 := reflect.TypeOf(obj)
	objName := t1.Name()
	objName = strings.ToLower(objName[:1]) + objName[1:]
	if prefix != "" {
		prefix = prefix + config2.Concat + objName
	} else {
		prefix = objName
	}

	for i, n := 0, t1.NumField(); i < n; i++ {
		field := t1.Field(i)
		name := field.Name
		fieldValue := v1.FieldByName(name)
		switch fieldValue.Kind() {
		case reflect.Struct:
			c.FillAttr(fieldValue.Interface(), prefix)
		default:
			name = strings.ToLower(name[:1]) + name[1:]
			val := getVal(prefix+config2.Concat+name, config2.ConfigMap)
			if _, ok := val.(*map[string]any); !ok {
				fieldValue.Elem().Set(reflect.ValueOf(val))
			}
		}
	}
}
