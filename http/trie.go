package http

import "strings"

const SEP = "/"

type Route struct {
	Method  string        `json:"method"`
	Path    string        `json:"path"`
	Handler HandlersChain `json:"handler"`
}

//trie 前缀树，用于存储路由服务和路径的映射关系
type trie struct {
	//next 下个节点信息
	next map[string]*trie
	//target 节点对应的路由规则
	target *Route

	key string

	isWord bool

	//IsInitial 是否是初始化
	IsInitial bool
}

//NewTrie 初始化树，返回一个前缀树的指针
func NewTrie() *trie {
	return &trie{
		next:      make(map[string]*trie),
		isWord:    false,
		IsInitial: true,
	}
}

//Insert 插入数据，其结构示例:
/*
- GET
  - /api
  - /test
    - /:id
      - /other -> Route{}
    - /test -> Route{}
*/
func (t *trie) Insert(word string, route *Route) {
	if !strings.Contains(word, "*") && !strings.Contains(word, ":") {
		newT := new(trie)
		newT.next = make(map[string]*trie)
		newT.isWord = false
		t.next[word] = &trie{
			isWord: false,
			target: route,
		}
	}
	items := strings.Split(word, SEP)
	for index, v := range items {
		if v == "" {
			continue
		}
		if t.next[v] == nil {
			node := new(trie)
			node.next = make(map[string]*trie)
			node.isWord = false
			node.key = v
			t.next[v] = node

		}

		if index == len(items)-1 {
			t.next[v].target = route
			t.next[v].isWord = true
			//终止循环
			break
		}
		t = t.next[v]
	}
}

//Search 前缀树查找，其时间复杂度为O(K),K=前缀长度，例如Path = /api/test/test,则前缀长度K=3
//注意：查询时需要注意路径参数匹配
func (t trie) Search(word string) *Route {
	if r, ok := t.next[word]; ok {
		return r.target
	}
	split := strings.Split(word, SEP)
	for idx, v := range split {
		if v == "" {
			continue
		}
		//if t.isWord {
		//	return t.target
		//}

		if t.next[v] == nil {
			//todo 路径参数匹配
			isPathVar := false
			for s, trie := range t.next {
				if strings.Contains(s, "*") {
					return trie.target

				}
				if strings.HasPrefix(s, ":") {
					t = *trie
					if idx == len(split)-1 {
						//已经是最后一个直接返回
						return t.target
					} else if len(t.next) == 0 {
						return t.target
					}
					isPathVar = true
					break
				}
			}
			if isPathVar {
				continue
			}
			return t.target
		}
		t = *t.next[v]
	}
	return nil
}
