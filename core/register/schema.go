package register


type reqRegister struct {
    WorkerID string `json:"wid"`
    Nonce    int64  `json:"nonce"`
    Version  string `json:"version"`
    Key      string `json:"key"`
}
// ChaosReader 读取服务器混淆后的数据
type ChaosReader struct {
    Bytes  []byte
    Offset int
}