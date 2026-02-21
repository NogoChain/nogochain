package types

import "io"

// Msg 网络消息
type Msg struct {
	Code    uint64
	Size    uint32
	Payload io.Reader
}
