package serve

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"

	log "jxcore/lowapi/logger"
)

var once sync.Once
var addrList []*net.IPNet = make([]*net.IPNet, 0)

func init() {
	once.Do(listNetInterfaces)
}

// APIResult API 返回结果
type APIResult struct {
	Data        *map[string]interface{} `json:"data,omitempty"`
	Description string                  `json:"desc"`
}

// NewAPIResult 使用对象
func NewAPIResult(data *map[string]interface{}) *APIResult {

	return &APIResult{Data: data, Description: "success"}
}

// WriteData 将结果写入相应
func WriteData(w http.ResponseWriter, data *map[string]interface{}) {
	WriteResult(w, NewAPIResult(data))
}

// WriteResult 返回结果
func WriteResult(w http.ResponseWriter, result *APIResult) {
	w.Header().Set("Content-Type", mimeJSON)
	w.WriteHeader(http.StatusOK)

	rs, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(rs)
}

// WriteSucess 写入标记操作为成功的空响应
func WriteSucess(w http.ResponseWriter) {
	WriteData(w, nil)
}

// listNetInterfaces 枚举网卡
func listNetInterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range interfaces {

		if (i.Flags&net.FlagLoopback) != 0 || (i.Flags&net.FlagPointToPoint) != 0 {
			continue
		}

		if addrs, err := i.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok {
					addrList = append(addrList, ipNet)
				}
			}
		}

	}

}

func pickIPNet(remoteAddr string) *net.IPNet {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil
	}
	ip := net.ParseIP(host)

	for _, it := range addrList {
		if it.Contains(ip) {
			return it
		}
	}
	return nil
}

// PickLocalAddr 根据远端地址，选择合适的 IP地址
func PickLocalAddr(r *http.Request) string {
	val := r.Context().Value(http.LocalAddrContextKey)
	if val != nil {
		addr := val.(net.Addr)
		if addr != nil {
			switch v := addr.(type) {
			case *net.TCPAddr:
				return v.IP.String()
			case *net.UDPAddr:
				return v.IP.String()
			case *net.IPAddr:
				return v.IP.String()
			}
		}
	}
	remoteAddr := r.RemoteAddr
	ipNet := pickIPNet(remoteAddr)
	if ipNet != nil {
		return ipNet.IP.String()
	}
	return "127.0.0.1"
}
