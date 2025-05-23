package dvap

import (
	"github.com/tidwall/gjson"
)

type SubModifyFunc func(p, s gjson.Result) (gjson.Result, gjson.Result)

// HasManyV2 准备废弃，请使用替代方法 NewDataer().HasMany
func HasManyV2(input interface{}, subGroup interface{}, relation string, f CompareFun, smf SubModifyFunc) (interface{}, error) {

	input_v := VtoJson(input)
	sub_g_v := VtoJson(subGroup)

	if !sub_g_v.IsArray() {
		// return input, fmt.Errorf("拆分数据必须是一个切片")
		sub_g_v = gjson.Parse("[]")
	}

	var result []interface{}

	if input_v.IsArray() {
		for _, iv := range input_v.Array() {

			// break
			// var remove []gjson.Result
			var filter []gjson.Result
			// FilterStructSlice(sub_g_v.Array(), &remove, &filter, func(sv gjson.Result) bool {
			// 	return f(iv, sv)
			// })

			for _, sv := range sub_g_v.Array() {
				if f(iv, sv) {
					if smf != nil {
						_iv, _sv := smf(iv, sv)
						iv = _iv
						filter = append(filter, _sv)
					} else {
						filter = append(filter, sv)
					}
				}
			}

			// fmt.Println(len(filter))
			result = append(result, VSetV(iv, JArrToInterface(filter), relation).Value())
		}

		var ret interface{} = result
		return ret, nil
	} else {
		// var remove []gjson.Result
		var filter []gjson.Result
		// FilterStructSlice(sub_g_v.Array(), &remove, &filter, func(sv gjson.Result) bool {
		// 	return f(input_v, sv)
		// })

		for _, sv := range sub_g_v.Array() {
			if f(input_v, sv) {
				if smf != nil {
					_input_v, _sv := smf(input_v, sv)
					input_v = _input_v
					filter = append(filter, _sv)
				} else {
					filter = append(filter, sv)
				}
			}
		}

		return VSetV(input_v, JArrToInterface(filter), relation).Value(), nil
	}

}

// HasOneV2 准备废弃，请使用替代方法 NewDataer().HasOne
func HasOneV2(input interface{}, subGroup interface{}, relation string, f CompareFun, smf SubModifyFunc) (interface{}, error) {
	input_v := VtoJson(input)
	sub_g_v := VtoJson(subGroup)

	if !sub_g_v.IsArray() {
		// return input, fmt.Errorf("拆分数据必须是一个切片")
		sub_g_v = gjson.Parse("[]")

	}

	var result []interface{}

	if input_v.IsArray() {
		for _, iv := range input_v.Array() {
			var match_v gjson.Result
			SliceFind(sub_g_v.Array(), &match_v, func(sv gjson.Result) bool {
				return f(iv, sv)
			})

			if smf != nil {
				iv, match_v = smf(iv, match_v)
			}
			result = append(result, VSetV(iv, match_v.Value(), relation).Value())
		}

		var ret interface{} = result
		return ret, nil
	} else {

		var match_v gjson.Result
		SliceFind(sub_g_v.Array(), &match_v, func(sv gjson.Result) bool {
			return f(input_v, sv)
		})
		if smf != nil {
			input_v, match_v = smf(input_v, match_v)
		}
		return VSetV(input_v, match_v.Value(), relation).Value(), nil

	}
}
