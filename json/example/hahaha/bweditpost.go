package bweditpost

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuchensheng/bintools/json/consts"
	"github.com/kuchensheng/bintools/json/executor/check"
	"github.com/kuchensheng/bintools/json/executor/js"
	"github.com/kuchensheng/bintools/json/executor/parameter"
	"github.com/kuchensheng/bintools/json/executor/predicate"
	"github.com/kuchensheng/bintools/json/executor/response"
	"github.com/kuchensheng/bintools/json/executor/server"
	"github.com/kuchensheng/bintools/json/executor/util"
	"github.com/kuchensheng/bintools/json/model"
	"github.com/kuchensheng/bintools/tracer/trace"
	"github.com/rs/zerolog/log"
)

var tenantId = func() string {
	var t string
	t = "hahaha"
	return t
}()

var steps = func() []model.ApixStep {
	var steps []model.ApixStep
	//将字符串内容初始化
	var strStep = `[{"prevId":"","graphId":"21330edbcb3242c1aacddfb69860d3fc","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"let res = {\n    \"code\":200,\n    \"message\":\"成功\",\n    \"data\":{}\n}\nres\n\n"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"21330edbcb3242c1aacddfb69860d3fc","graphId":"8ad2dea3b33c4f4d9f2b82102c1016e0","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":[{"enabled":true,"type":"","key":"$req.data.plateNumber","value":"浙","operator":"inc","isRegex":false,"cases":null}],"predicateType":0,"thenGraphId":"23f398a93f764236b0c2b7a05d6a1feb","elseGraphId":"9c719f63a18a459b81e54e8cb795f9d5"},{"prevId":"8ad2dea3b33c4f4d9f2b82102c1016e0","graphId":"9c719f63a18a459b81e54e8cb795f9d5","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"8ad2dea3b33c4f4d9f2b82102c1016e0","graphId":"23f398a93f764236b0c2b7a05d6a1feb","code":"","domain":"10.30.30.95:38800","protocol":"http","method":"GET","path":"/api/tddm/model/basic/1600380685216194561","parameters":[{"name":"token","type":"String","in":"header","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"default":"$req.header.token","required":true},{"name":"personTag","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"personName","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"personMobile","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"plateNumber","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"authStartTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"authEndTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"valid","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"remarks","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"factoryIdentify","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"factoryRecordIdentify","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"createTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"updateTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"id","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"default":"$req.data.bwId","required":false},{"name":"pageSize","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":true},{"name":"pageNum","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":true},{"name":"root","type":"","in":"body","schema":{"type":"object","properties":{},"subtype":"","children":null,"default":""},"required":false}],"local":false,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"23f398a93f764236b0c2b7a05d6a1feb","graphId":"f909fb3bfd9a429ebe55c6b3fbe37e32","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"let records = $23f398a93f764236b0c2b7a05d6a1feb.$resp.data.data.records\nif(records.length === 0){\n    null\n}else{\n    let ok = {\n        \"blackGuid\":records[0].factoryRecordIdentify,\n        \"factoryIdentify\":records[0].factoryIdentify\n    }\n    ok\n}"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"f909fb3bfd9a429ebe55c6b3fbe37e32","graphId":"d782b51c2a6f4132a8dc17edcef4f08f","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":[{"enabled":true,"type":"","key":"$f909fb3bfd9a429ebe55c6b3fbe37e32.$resp.export","value":"","operator":"nil","isRegex":false,"cases":null}],"predicateType":0,"thenGraphId":"f71d93723d40403f9f1d823570e93d55","elseGraphId":"d5f586390e944368a89668eb5cd60ac1"},{"prevId":"d782b51c2a6f4132a8dc17edcef4f08f","graphId":"d5f586390e944368a89668eb5cd60ac1","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"d782b51c2a6f4132a8dc17edcef4f08f","graphId":"f71d93723d40403f9f1d823570e93d55","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export = {\n    \"code\":100,\n    \"message\":\"记录不存在\"\n}\n\n"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"d5f586390e944368a89668eb5cd60ac1","graphId":"53f61088553443fdb0225de6d3637785","code":"","domain":"isc-orchestration-service:38234","protocol":"http","method":"POST","path":"/api/app/orc/integration/bw/edit","parameters":[{"name":"root","type":"","in":"body","schema":{"type":"object","properties":{"blackGuid":{"name":"blackGuid","default":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export.blackGuid","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"blackWhiteType":{"name":"blackWhiteType","default":"2","in":"body","type":"int","subtype":"","children":null,"properties":null,"required":false},"endDate":{"name":"endDate","default":"$req.data.authEndTime","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"factoryRecordIdentify":{"name":"factoryRecordIdentify","default":"$f909fb3bfd9a429ebe55c6b3fbe37e32.$resp.export.factoryIdentify","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"plateNumber":{"name":"plateNumber","default":"$req.data.plateNumber","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"startDate":{"name":"startDate","default":"$req.data.authStartTime","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false}},"subtype":"","children":null,"default":""},"required":false}],"local":false,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"53f61088553443fdb0225de6d3637785","graphId":"17ecfcd454ef4bbcaac7d66575a5d7e5","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":[{"enabled":true,"type":"","key":"$53f61088553443fdb0225de6d3637785.$resp.data.code","value":"200","operator":"==","isRegex":false,"cases":null}],"predicateType":0,"thenGraphId":"5774760d68f74722855e236a4a66647f","elseGraphId":"4f26cb6ad3714729bbe946b7702ef623"},{"prevId":"17ecfcd454ef4bbcaac7d66575a5d7e5","graphId":"4f26cb6ad3714729bbe946b7702ef623","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"17ecfcd454ef4bbcaac7d66575a5d7e5","graphId":"5774760d68f74722855e236a4a66647f","code":"","domain":"10.30.30.95:38800","protocol":"http","method":"PATCH","path":"/api/tddm/model/basic/1600380685216194561","parameters":[{"name":"token","type":"String","in":"header","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"default":"$req.header.token","required":true},{"name":"personTag","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"personName","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"personMobile","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"plateNumber","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"authStartTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"authEndTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"valid","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"remarks","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"factoryIdentify","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"factoryRecordIdentify","type":"String","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"createTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"updateTime","type":"Timestamp","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"required":false},{"name":"id","type":"Integer","in":"query","schema":{"type":"","properties":null,"subtype":"","children":null,"default":""},"default":"$req.data.bwId","required":false},{"name":"root","type":"","in":"body","schema":{"type":"object","properties":{"authEndTime":{"name":"authEndTime","default":"$req.data.authEndTime","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"authStartTime":{"name":"authStartTime","default":"$req.data.authStartTime","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"createTime":{"name":"createTime","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"factoryIdentify":{"name":"factoryIdentify","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"factoryRecordIdentify":{"name":"factoryRecordIdentify","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"id":{"name":"id","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"personMobile":{"name":"personMobile","default":"$req.data.personMobile","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"personName":{"name":"personName","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"personTag":{"name":"personTag","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"plateNumber":{"name":"plateNumber","default":"$req.data.plateNumber","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"remarks":{"name":"remarks","default":"$req.data.remarks","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"updateTime":{"name":"updateTime","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"valid":{"name":"valid","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false}},"subtype":"","children":null,"default":""},"required":false}],"local":false,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"5774760d68f74722855e236a4a66647f","graphId":"fc7897405550431db34dac1cb5885031","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":[{"enabled":true,"type":"","key":"$5774760d68f74722855e236a4a66647f.$resp.data.code","value":"0","operator":"==","isRegex":false,"cases":null}],"predicateType":0,"thenGraphId":"c83199232f3e468ea737564d74326576","elseGraphId":"fcdfe4b1c39046d3896475d6c3fca64f"},{"prevId":"fc7897405550431db34dac1cb5885031","graphId":"fcdfe4b1c39046d3896475d6c3fca64f","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"","script":{"language":"","script":""},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"fc7897405550431db34dac1cb5885031","graphId":"c83199232f3e468ea737564d74326576","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export = {\n    \"code\":200,\n    \"message\":\"编辑成功\"\n}\n\n"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"fcdfe4b1c39046d3896475d6c3fca64f","graphId":"a5100e25c0634080827f893010791afd","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export = {\n    \"code\":$5774760d68f74722855e236a4a66647f.$resp.data.code,\n    \"message\": $5774760d68f74722855e236a4a66647f.$resp.data.message\n}"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"4f26cb6ad3714729bbe946b7702ef623","graphId":"8a3a80ced5d348f88404a597be5e9699","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export = {\n    \"code\":$53f61088553443fdb0225de6d3637785.$resp.data.code,\n    \"message\":$53f61088553443fdb0225de6d3637785.$resp.data.message\n}\n\n"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""},{"prevId":"9c719f63a18a459b81e54e8cb795f9d5","graphId":"20a315cf2b9548fb9a6870115a48c365","code":"","domain":"","protocol":"","method":"","path":"","parameters":null,"local":true,"language":"javascript","script":{"language":"javascript","script":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export = {\n    \"code\":100,\n    \"message\":\"车牌号不合法\"\n}\n\n"},"predicate":null,"predicateType":0,"thenGraphId":"","elseGraphId":""}]`
	if err := json.Unmarshal([]byte(strStep), &steps); err != nil {
		log.Error().Msgf("步骤信息初始化失败,%v", err)
	}
	return steps
}()
var parameters = func() []model.ApixParameter {
	var parameters []model.ApixParameter
	//将字符串内容初始化
	var strParameter = `[{"name":"root","type":"","in":"body","schema":{"type":"object","properties":{"authEndTime":{"name":"authEndTime","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"authStartTime":{"name":"authStartTime","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"bwId":{"name":"bwId","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"personMobile":{"name":"personMobile","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"plateNumber":{"name":"plateNumber","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false},"remarks":{"name":"remarks","default":"","in":"body","type":"string","subtype":"","children":null,"properties":null,"required":false}},"subtype":"","children":null,"default":""},"required":false}]`
	if err := json.Unmarshal([]byte(strParameter), &parameters); err != nil {
		log.Error().Msgf("参数初始化失败,%v", err)
	}
	return parameters
}()

var stepMaps = listToMap(steps)

var responses = func() map[string]model.ApixResponse {
	responseMap := make(map[string]model.ApixResponse)
	//字符串初始化
	var strResponse = `{"200":{"schema":{"type":"object","properties":{"code":{"name":"code","default":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export.code","in":"","type":"string","subtype":"","children":null,"properties":null,"required":false},"message":{"name":"message","default":"$21330edbcb3242c1aacddfb69860d3fc.$resp.export.message","in":"","type":"string","subtype":"","children":null,"properties":null,"required":false}},"subtype":"","children":null,"default":""},"setCookie":null}}`
	if err := json.Unmarshal([]byte(strResponse), &responseMap); err != nil {
		log.Error().Msgf("响应信息初始化失败,%v", err)
	}
	return responseMap
}()

func Executorbweditpost(ctx *gin.Context) (any, error) {
	defer func() {
		if x := recover(); x != nil {
			log.Error().Msgf("请求执行失败,%v", x)
		}
	}()
	log.Info().Msgf("当前请求:%s,Method:%s,执行文件:[bweditpost]", ctx.Request.URL.Path, ctx.Request.Method)
	if err := check.CheckTenantId(ctx, tenantId); err != nil {
		return nil, err
	}
	parameterMap := parameter.SetParameterMap(ctx)
	log.Info().Msg("链路跟踪启动...")
	tracer := trace.NewServerTracer(ctx.Request)
	log.Info().Msgf("traceId = %s", tracer.TracId)

	//必填参数检查
	if err := parameter.CheckParameter(parameters, parameterMap); err != nil {
		log.Warn().Msgf("缺少必填参数:%v", err)
		tracer.EndServerTracer(trace.WARNING, err.Error())
		return nil, err
	}
	ctx.Set(consts.RESULTMAP, make(map[string]any))
	ctx.Set(consts.TRACER, tracer)
	//执行步骤
	log.Info().Msgf("参数校验通过,开始执行逻辑流程")
	if err := executeStep(ctx, "", steps); err != nil {
		log.Warn().Msgf("流程执行失败:%v", err)
		tracer.EndServerTracer(trace.WARNING, err.Error())
		return nil, err
	}
	log.Info().Msgf("流程步骤执行完毕，开始组装结果映射...")
	return response.BuildSuccessResponse(ctx, responses)

}

func executeAll(ctx *gin.Context) error {
	return nil
}

//todo 下个迭代，这里将直接用模板生成执行代码，而不是先解析再执行
//executeStep 执行步骤
func executeStep(ctx *gin.Context, PrevId string, sts []model.ApixStep) error {
	defer deferHandler()
	subList := func(parentId string) []model.ApixStep {
		var result []model.ApixStep
		for _, st := range sts {
			if st.PrevId == parentId {
				result = append(result, st)
			}
		}
		return result
	}

	subSts := subList(PrevId)
	if len(subSts) < 1 {
		return nil
	}
	for _, step := range subSts {
		if err := runStep(step, ctx, stepMaps); err != nil {
			return err
		}
		//执行子节点
		if err := executeStep(ctx, step.GraphId, sts); err != nil {
			return err
		}
	}
	return nil
}

func listToMap(steps []model.ApixStep) map[string]model.ApixStep {
	result := make(map[string]model.ApixStep)
	for _, step := range steps {
		result[step.GraphId] = step
	}
	return result
}

func runStep(step model.ApixStep, ctx *gin.Context, stepMap map[string]model.ApixStep) error {
	log.Info().Msgf("执行步骤节点:%s", step.GraphId)

	if step.Language == "javascript" {
		// 执行JS脚本内容
		if result, e := js.ExecuteJavaScript(ctx, step.Script.Script, step.GraphId); e != nil {
			return e
		} else {
			util.SetResultValue(ctx, fmt.Sprintf("%s%s%s", consts.KEY_TOKEN, step.GraphId, ".$resp.export"), result)
		}
	} else if step.Predicate != nil {
		//执行判断逻辑
		if ok, e := predicate.ExecPredicates(ctx, step.Predicate, step.PredicateType); e != nil {
			return e
		} else {
			nextStep := stepMap[step.ThenGraphId]
			if !ok {
				nextStep = stepMap[step.ElseGraphId]
			}
			if e = runStep(nextStep, ctx, stepMap); e != nil {
				return e
			}
		}
	} else {
		//执行普通的服务请求
		if e := server.ExecServer(ctx, step); e != nil {
			return e
		}
	}
	return nil
}

func deferHandler() error {
	if x := recover(); x != nil {
		return x.(error)
	}
	return nil
}
