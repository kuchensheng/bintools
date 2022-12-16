package js

import (
	"strings"
	"testing"
)

func TestExecuteJavaScript(t *testing.T) {
	var script = "let records = $23f398a93f764236b0c2b7a05d6a1feb.$resp.data.data.records\nif(records.length === 0){\n    null\n}else{\n    let ok = {\n        \"blackGuid\":records[0].factoryRecordIdentify,\n        \"factoryIdentify\":records[0].factoryIdentify\n    }\n    ok\n}"

	split := strings.Split(script, "\n")

	for _, s := range split {
		println(s)
	}
}
