package consts

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"strings"
	"time"
)

var Cache = cache.New(time.Minute, 30*time.Second)

var GlobalPrefix = "/api/app/orc/"
var GlobalTestPrefix = "/api/app/test/orc/"

const GlobalTemplate = "tmp.tmpl"

func init() {
	out := zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: true,
	}
	out.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%-4s]", i))
	}
	out.FormatCaller = func(i interface{}) string {
		_, f, line, _ := runtime.Caller(7)
		return fmt.Sprintf("%s:%2d:", f, line)
	}
	out.TimeFormat = time.RFC3339Nano //"2006-01-02 15:04:05"
	out.FormatTimestamp = func(i interface{}) string {
		now := time.Now()
		return fmt.Sprintf("%00s,%d", now.Format("2006-01-02 15:04:05"), now.Nanosecond()/1e6)
	}
	log.Logger = log.Logger.Output(out).With().Caller().Logger()
}

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
