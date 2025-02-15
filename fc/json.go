package fc

import (
	sj "github.com/guyannanfei25/go-simplejson"
)

//多个数组json处理
func JsonDecodes(js []*sj.Json) []string {
	arr := make([]string, 0, len(js))
	for _, j := range js {
		v := JsonDecode(j)
		if v != "" {
			arr = append(arr, v)
		}
	}
	return arr
}

//用于查看对应json结果的字符串表示
func JsonDecode(j *sj.Json) string {
	if j == nil {
		return ""
	}
	bytes, err := j.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

//返回第一个元素中没有的数据
func JsonMerge(j1, j2 *sj.Json) *sj.Json {
	j2m, err := j2.Map()
	if err != nil {
		panic(err)
	}
	for key, val := range j2m {
		j1.Set(key, val)
	}
	return j1
}

//获取json对应的实际大小
func JsonBytesLength(jsons []*sj.Json) int64 {
	var b int64
	for _, j := range jsons {
		if j == nil {
			continue
		}
		jencode, err := j.Encode()
		if err != nil {
			panic(err)
		}
		b += int64(len(jencode))
	}
	return b
}

//特殊逻辑, 用于检测proc_chain是否传入parent_md5
func IsSetQueryFieldParentMd5(query *sj.Json) bool {
	queryField, ok := query.Get("post").CheckGet("query_field")
	if !ok {
		return false
	}
	field := queryField.MustString()
	if field == "parent_md5" {
		return true
	}
	return false
}
