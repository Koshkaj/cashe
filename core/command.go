package core

import (
	"bytes"
	"encoding/binary"
)

type Command byte

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

type CommandGet struct {
	Key []byte
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdGet)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdSet)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	valueLen := int32(len(c.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, c.Value)
	binary.Write(buf, binary.LittleEndian, int32(c.TTL))
	return buf.Bytes()
}

type CommandJoin struct {
	RaftAddr []byte
	NodeID   []byte
}

func (c *CommandJoin) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdJoin)
	var (
		raftAddrLen = int32(len(c.RaftAddr))
		nodeIDLen   = int32(len(c.NodeID))
	)
	binary.Write(buf, binary.LittleEndian, raftAddrLen)
	binary.Write(buf, binary.LittleEndian, c.RaftAddr)

	binary.Write(buf, binary.LittleEndian, nodeIDLen)
	binary.Write(buf, binary.LittleEndian, c.NodeID)
	return buf.Bytes()

}

type CommandDel struct {
	Key []byte
}

func (c *CommandDel) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdDel)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}
