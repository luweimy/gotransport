package codec

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

type ProtobufCodec struct {
}

func (c ProtobufCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Message); ok {
		return proto.Marshal(m)
	}

	return nil, fmt.Errorf("%T is not a proto.Marshaler", i)
}

func (c ProtobufCodec) Decode(data []byte, i interface{}) error {
	if m, ok := i.(proto.Message); ok {
		return proto.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a proto.Unmarshaler", i)
}
