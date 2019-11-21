package utils

import "reflect"

//Hash Hash
func Hash(a interface{}, b interface{}) []string {
	set := make([]string, 0)
	hash := make(map[interface{}]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).String()
		if _, found := hash[el]; found {
			set = append(set, el)
		}
	}

	return set
}
