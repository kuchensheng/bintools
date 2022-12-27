package consts

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = cache.New(time.Minute, 30*time.Second)

var GlobalPrefix = "/api/app/orc/"

const GlobalTemplate = "tmp.tmpl"

//func init() {
//	//删除文件
//	go func() {
//		ticker := time.NewTicker(time.Minute)
//		for {
//			select {
//			case <-ticker.C:
//				getwd, _ := os.Getwd()
//
//				if globs, err := ioutil.ReadDir(filepath.Join(getwd, "scripts")); err != nil {
//					log.Error().Msgf("获取js文件列表失败,%v", err)
//				} else {
//					for _, glob := range globs {
//						key := strings.ReplaceAll(glob.Name(), ".js", "")
//						if _, ok := Cache.Get(key); !ok {
//							fp := filepath.Join(getwd, "scripts", glob.Name())
//							log.Info().Msgf("删除文件：%s", fp)
//							_ = os.Remove(fp)
//						}
//					}
//				}
//			case <-time.After(time.Second * 30):
//				continue
//			}
//		}
//	}()
//}

const (
	KEY_TOKEN         = "$"
	KEY_REQ_CONNECTOR = "."
	KEY_REQ           = "$req"
	KEY_REQ_BODY      = KEY_REQ + KEY_REQ_CONNECTOR + "data"
	KEY_BODY          = "body"
	KEY_QUERY         = "query"
	KEY_HEADER        = "header"
	KEY_COOKIE        = "cookie"
	KEY_PATH          = "path"
	KEY_FORM          = "form"
	KEY_DATA          = "data"
	KEY_REQ_QUERY     = KEY_REQ + KEY_REQ_CONNECTOR + KEY_QUERY
	KEY_REQ_HEADER    = KEY_REQ + KEY_REQ_CONNECTOR + KEY_HEADER
	KEY_REQ_COOKIE    = KEY_REQ + KEY_REQ_CONNECTOR + KEY_COOKIE
	KEY_REQ_PATH      = KEY_REQ + KEY_REQ_CONNECTOR + KEY_PATH
	KEY_REQ_FORM      = KEY_REQ + KEY_REQ_CONNECTOR + KEY_FORM
)

const (
	RESULTMAP    = "resultMap"
	PARAMETERMAP = "parameterMap"

	TRACER = "tracer"
)

const (
	OBJECT = "object"
	ARRAY  = "array"
)

const (
	TENANT_ID = "isc-tenant-id"
)
