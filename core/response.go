package core

import (
	"bytes"
	"encoding/binary"
)

type ResponseSet struct {
	Status Status
}

func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
}

type ResponseDel struct {
	Status Status
}

func (r *ResponseDel) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)

	valueLen := int32(len(r.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, r.Value)
	return buf.Bytes()
}

type ResponseHas struct {
	Status Status
	Value  bool
}

func (r *ResponseHas) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	binary.Write(buf, binary.LittleEndian, r.Value)
	return buf.Bytes()
}
