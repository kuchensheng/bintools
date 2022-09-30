package model

type ApixProperty struct {
	Name       string                  `json:"name"`       // 属性名称
	Default    string                  `json:"default"`    // 属性的默认值
	In         string                  `json:"in"`         // 属性的填充位置，resBody，只有在响应参数内才看这个值
	Type       string                  `json:"type"`       // 属性的数据类型，如果是 array，则需要读 subType
	SubType    string                  `json:"subtype"`    // 属性的子类型，如果是 object，就需要读取 Children 的内容
	Children   []ApixProperty          `json:"children"`   // 数组属性所对应的子属性，属性类型为 array 并且 subType 为 object 时，读取这里的子属性
	Properties map[string]ApixProperty `json:"properties"` // 属性类型为 object 时，其对应的子对象所具备的属性
	Required   bool                    `json:"required"`   // 是否必填
}

type ApixSchema struct {
	Type       string                  `json:"type"`       // Schema 的类型，object 或 array
	Properties map[string]ApixProperty `json:"properties"` // Schema 对象所具备的属性
	SubType    string                  `json:"subtype"`    // 属性的子类型，如果是 object，就需要读取 Children 的内容
	Children   []ApixProperty          `json:"children"`   // 数组属性所对应的子属性，属性类型为 array 并且 subType 为 object 时，读取这里的子属性
	Default    string                  `json:"default"`    // 属性的默认值, Type 不是 object 或 array 时，将取得表达式所描述的值并予以返回，Type 是 object 时，将获取表达式对应的对象或数组，并直接返回
}

type ApixParameter struct {
	Name     string     `json:"name"`              // 参数名称
	Type     string     `json:"type"`              // 参数数据类型
	In       string     `json:"in"`                // 参数的填充位置，query/body/path，为path的时候只能放在最后
	Schema   ApixSchema `json:"schema"`            // 参数的具体信息，当 type 为 object 或 array 时生效
	Default  string     `json:"default,omitempty"` // 参数的默认值，可以从全局 request 或上一个节点的 response 内获取数据
	Required bool       `json:"required"`          // 是否必填
}

type ApixApi struct {
	Path         string          `json:"path"`         // Api 的请求路径，该路径在OS内全局唯一
	Protocol     string          `json:"protocol"`     // 请求协议，只能是http或https
	Method       string          `json:"method"`       // 请求方法，GET/POST/PUT/DELETE/OPTION
	Domain       string          `json:"domain"`       // Api 的域名（可选，不填的情况默认refer为当前服务）
	Parameters   []ApixParameter `json:"parameters"`   // Api 的请求参数
	RequireLogin bool            `json:"requireLogin"` // API 是否需要登录
	Version      string          `json:"version,omitempty"`
}

type ApixSetCookie struct {
	Name    string `json:"name"`    // 要写入的 Cookie 名称
	Default string `json:"default"` //  要设置的值的来源
}

type ApixResponse struct {
	Schema    ApixSchema               `json:"schema"`    // Api 响应的数据结构定义
	SetCookie map[string]ApixSetCookie `json:"setCookie"` // Api 响应后，要写入的 cookie 内容
}

type ApixScript struct {
	Language string `json:"language"` // 脚本语言的类型，目前只能填 javascript/goscript
	Script   string `json:"script"`   // Api 替换为脚本代码后，所需要执行的脚本代码
}

type ApixSwitchPredicate struct {
	IsDefault   bool   `json:"isDefault"`   // 是否 switch 的默认节点，即全部的条件都不命中，所触发的流程
	Key         string `json:"key"`         // 如果父节点有 Key，则此处的 key 无效
	Value       string `json:"value"`       // 要判断的值
	Operator    string `json:"operator"`    // 判断操作符，当节点型有 Key 时，操作符永远为 ==(等于)，否则这里的操作符逻辑将等同于父节点的 Operator 逻辑
	ThenGraphId string `json:"thenGraphId"` // case 条件命中时执行的节点
}

type ApiStepPredicate struct {
	Enabled  bool                  `json:"enabled"`  // 是否启用条件判断
	Type     string                `json:"type"`     // 条件的类型，只能是 if/switch
	Key      string                `json:"key"`      // 要判断的名
	Value    string                `json:"value"`    // 要判断的值
	Operator string                `json:"operator"` // 判断操作符，可以是以下的符号：==(等于),!=(不等于),>(大于),>=(大于等于),<(小于),<=(小于等于),in(在集合内),!in(不在集合内),inc(包含),!inc(不包含),nil(空),!nil(非空),true(真),false(假)
	IsRegex  bool                  `json:"isRegex"`  // 当此选项为true时，value为正则表达式，将判断key的值是否与value所指的正则表达式匹配，此时Operator会失效
	Cases    []ApixSwitchPredicate `json:"cases"`    // 针对 type 是 switch 的情况才有效，指出switch所属的case分支，key有值的情况下，对key进行switch，key没值的情况，转换为when语句，这里的switch没有落入机制
}

type ApixStep struct {
	PrevId     string          `json:"prevId"`     // 上一个节点的 GraphId，没有上一个节点的情况，此项为0
	GraphId    string          `json:"graphId"`    // 节点的Id，在一个ApixData内，此Id唯一
	Code       string          `json:"code"`       // StepId （编译器不处理，但是保留）
	Domain     string          `json:"domain"`     // 请求 Api 时的域名和端口，如 isc-permission-service:32100
	Protocol   string          `json:"protocol"`   // 请求协议，只能是http或https
	Method     string          `json:"method"`     // 请求方法，GET/POST/PUT/DELETE/OPTION
	Path       string          `json:"path"`       // 请求路径
	Parameters []ApixParameter `json:"parameters"` // 请求的参数
	//ApiId         string                 `json:"apiId"`         // 数据库中的ApiId（编译器不处理，但是保留）
	Local         bool               `json:"local"`         // 是否本地代码，如果要使用Script来代替这个节点的请求，则置为true
	Language      string             `json:"language"`      // 脚本语言的类型，目前只能填 javascript/goscript
	Script        ApixScript         `json:"script"`        // 脚本代码
	Predicate     []ApiStepPredicate `json:"predicate"`     // 条件判断，如果是switch，则只能有一个节点
	PredicateType int                `json:"predicateType"` // 条件判断模式，0:所有条件都为真，1:任意条件为真，只有一个条件时不生效
	ThenGraphId   string             `json:"thenGraphId"`   // predicate 条件命中时执行的节点
	ElseGraphId   string             `json:"elseGraphId"`   // predicate 条件不命中时执行的节点
}

type ApixRule struct {
	Api      ApixApi                 `json:"api"`       // 最终可调用的 restful api 的定义
	Response map[string]ApixResponse `json:"responses"` // 最终可调用的 restful api 的返回数据
	Steps    []ApixStep              `json:"steps"`     // api 编排后的的执行步骤
}

// ApixData API的完整定义
type ApixData struct {
	Rule ApixRule `json:"rule"`
}
