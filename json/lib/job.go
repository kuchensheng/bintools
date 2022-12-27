package lib

import (
	"github.com/rs/zerolog/log"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var Wd, _ = os.Getwd()
var GoPath = "D:\\worksapace\\go"

var programMap = make(map[string]*interp.Interpreter)

var lock sync.Mutex

func PutProgramMap(key string, v *interp.Interpreter) {
	lock.Lock()
	defer lock.Unlock()
	programMap[key] = v
}

func GetProgramMap(key string) (v *interp.Interpreter, ok bool) {
	lock.Lock()
	defer lock.Unlock()
	v, ok = programMap[key]
	return
}

var ScriptEngineFunc = func() *interp.Interpreter {
	i := interp.New(interp.Options{GoPath: GoPath})
	i.Use(stdlib.Symbols)
	i.Use(Symbols)
	return i
}

func init() {

	//启动时预编译所有的json文件
	filepath.Walk(filepath.Join(Wd, "example"), func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			log.Info().Msgf("预加载文件:%s", path)
			var scriptEngine = ScriptEngineFunc()
			if _, err = scriptEngine.EvalPath(path); err != nil {
				log.Error().Msgf("Go文件无法被编译，%v", err)
				return err
			} else {
				log.Info().Msgf("文件[%s]加载完成", path)
				key := strings.ReplaceAll(info.Name(), ".go", "")
				PutProgramMap(key, scriptEngine)
			}
		}
		return nil
	})
}
