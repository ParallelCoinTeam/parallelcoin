package codec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	proto "github.com/gogo/protobuf/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/ugorji/go/codec"

	"github.com/apache/thrift/lib/go/thrift"
)

// Codec defines the interface that decode/encode payload.
type Codec interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

// ByteCodec uses raw slice pf bytes and don't encode/decode.
type ByteCodec struct{}

// Encode returns raw slice of bytes.
func (c ByteCodec) Encode(i interface{}) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}
	if data, ok := i.(*[]byte); ok {
		return *data, nil
	}

	return nil, fmt.Errorf("%T is not a []byte", i)
}

// Decode returns raw slice of bytes.
func (c ByteCodec) Decode(data []byte, i interface{}) error {
	reflect.Indirect(reflect.ValueOf(i)).SetBytes(data)
	return nil
}

// JSONCodec uses json marshaler and unmarshaler.
type JSONCodec struct{}

// Encode encodes an object into slice of bytes.
func (c JSONCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

// Decode decodes an object from slice of bytes.
func (c JSONCodec) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

// PBCodec uses protobuf marshaler and unmarshaler.
type PBCodec struct{}

// Encode encodes an object into slice of bytes.
func (c PBCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Marshaler); ok {
		return m.Marshal()
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Marshal(m)
	}

	return nil, fmt.Errorf("%T is not a proto.Marshaler", i)
}

// Decode decodes an object from slice of bytes.
func (c PBCodec) Decode(data []byte, i interface{}) error {
	if m, ok := i.(proto.Unmarshaler); ok {
		return m.Unmarshal(data)
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a proto.Unmarshaler", i)
}

// MsgpackCodec uses messagepack marshaler and unmarshaler.
type MsgpackCodec struct {
	handle  codec.MsgpackHandle
	bufsize int
	buf     *bytes.Buffer
	enc     *codec.Encoder
	dec     *codec.Decoder
}

// NewMsgpackCodec uses the optimised github.com/ugorji/go and we want to
// eliminate runtime allocations so a configurable buffer size can be set
// with the parameter and the codec can be reused
func NewMMsgpackCodec(bufsize int) *MsgpackCodec {
	cdc := MsgpackCodec{
		bufsize: bufsize,
		buf:     bytes.NewBuffer(make([]byte, bufsize)),
	}
	cdc.enc = codec.NewEncoder(cdc.buf, &cdc.handle)
	cdc.dec = codec.NewDecoder(cdc.buf, &cdc.handle)
	return &cdc
}

// Encode encodes an object into slice of bytes.
func (c *MsgpackCodec) Encode(i interface{}) (out []byte, e error) {
	c.buf.Reset()
	c.enc.Reset(c.buf)
	e = c.enc.Encode(i)
	if e != nil {
		return
	}
	out = c.buf.Bytes()
	return
}

// Decode decodes an object from slice of bytes.
func (c *MsgpackCodec) Decode(data []byte, i interface{}) error {
	if len(data) > c.bufsize {
		panic("buffer is too small for received data")
	}
	c.buf.Reset()
	c.dec.Reset(c.buf)
	c.buf.Read(data)
	return c.dec.Decode(i)
}

type ThriftCodec struct{}

func (c ThriftCodec) Encode(i interface{}) ([]byte, error) {
	b := thrift.NewTMemoryBufferLen(1024)
	p := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(b)
	t := &thrift.TSerializer{
		Transport: b,
		Protocol:  p,
	}
	t.Transport.Close()
	return t.Write(context.Background(), i.(thrift.TStruct))
}

func (c ThriftCodec) Decode(data []byte, i interface{}) error {
	t := thrift.NewTMemoryBufferLen(1024)
	p := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(t)
	d := &thrift.TDeserializer{
		Transport: t,
		Protocol:  p,
	}
	d.Transport.Close()
	return d.Read(i.(thrift.TStruct), data)
}
