package model

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
)

//BufioRead is
func BufioRead(name string) io.Reader {
	var reader io.Reader
	if fileObj, err := os.Open(name); err == nil {
		defer fileObj.Close()

		reader := bufio.NewReader(fileObj)

		if result, err := reader.ReadString(byte('@')); err == nil {
			fmt.Println("使用ReadSlince相关方法读取内容:", result)
		}

		return reader
	}
	return reader
}
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