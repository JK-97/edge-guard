package register

import "io"

// chaosReader 读取服务器混淆后的数据
type chaosReader struct {
	Bytes  []byte
	Offset int
}

func (r *chaosReader) read(p []byte) (n int, err error) {
	length := len(r.Bytes)
	remain := length - r.Offset
	if remain <= 0 {
		return 0, io.EOF
	}
	length = len(p)
	if length > remain {
		err = io.EOF
	} else {
		remain = length
	}

	for n = 0; n < remain; n++ {
		b := r.Bytes[r.Offset+n]
		if b >= 0x80 {
			p[n] = b - 0x80
		} else {
			p[n] = b + 0x80
		}
	}
	r.Offset += remain
	return
}
