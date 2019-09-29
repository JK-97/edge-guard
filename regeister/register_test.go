package regeister

import (
	"fmt"
	"testing"
)

type test struct {
	Data []string `json:"data"`
}

func TestIsInit(t *testing.T) {
	//data :="'{ 'data':['aaa','aaa']}'"
	//file, err := ioutil.ReadAll(strings.NewReader(data))
	//if err != nil {
	//	log.Info(err)
	//}
	//datalist := test{}
	//json.Unmarshal(file, &datalist)
	//fmt.Println(datalist)

	str := "#nameserver 192.168.0.66"
	//rip := regexp.MustCompile(" {0,}nameserver {1,}((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
	//res := rip.FindAllStringSubmatch(str, -1)
	//for _,per :=range res{
	//	fmt.Println(per[1])
	//}
	fmt.Println(string(str[0]))
}
