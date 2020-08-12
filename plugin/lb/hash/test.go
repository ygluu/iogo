package hash

import (
	"iogo/log"
)

func Test() {
	hash := NewLoadBalan(150)
	hash.Add("aaa", 0)
	hash.Add("bbb", 0)
	hash.Add("ccc", 0)
	countMap := make(map[string]int)
	countMap["aaa"] = 0
	countMap["bbb"] = 0
	countMap["ccc"] = 0

	for i := 0; i < 1000; i++ {
		// 一致性
		//v := hash.Get("dddddd")
		// 均衡
		v := hash.Get("")
		count := countMap[v]
		countMap[v] = count + 1
	}

	for k, v := range countMap {
		log.I(999, "%s:%d", k, v)
	}
}
