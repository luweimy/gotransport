package codec

import (
	"github.com/vmihailenco/msgpack"
)

type MsgPackCodec struct {
}

func (c MsgPackCodec) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (c MsgPackCodec) Decode(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
