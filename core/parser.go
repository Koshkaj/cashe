package core

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Status byte

func (s Status) String() string {
	switch s {
	case StatusError:
		return "ERR"
	case StatusOK:
		return "OK"
	case StatusKeyNotFound:
		return "KEYNOTFOUND"
	default:
		return "NONE"
	}
}

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}
	err := binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp, err
}

func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)

	resp.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &resp.Value)

	return resp, nil

}

func ParseDelResponse(r io.Reader) (*ResponseDel, error) {
	resp := &ResponseDel{}
	err := binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp, err
}

func ParseHasResponse(r io.Reader) (*ResponseHas, error) {
	resp := &ResponseHas{}
	binary.Read(r, binary.LittleEndian, &resp.Status)
	err := binary.Read(r, binary.LittleEndian, &resp.Value)
	return resp, err
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}
	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	case CmdJoin:
		return parseJoinCommand(r), nil
	case CmdLeave:
		return parseLeaveCommand(r), nil
	case CmdDel:
		return parseDelCommand(r), nil
	case CmdHas:
		return parseHasCommand(r), nil
	default:
		return nil, fmt.Errorf("invalid command %s", string(cmd))
	}
}

func parseJoinCommand(r io.Reader) *CommandJoin {
	cmd := &CommandJoin{}
	var (
		raftAddrLen int32
		nodeIDLen   int32
	)
	binary.Read(r, binary.LittleEndian, &raftAddrLen)
	cmd.RaftAddr = make([]byte, raftAddrLen)
	binary.Read(r, binary.LittleEndian, &cmd.RaftAddr)

	binary.Read(r, binary.LittleEndian, &nodeIDLen)
	cmd.NodeID = make([]byte, nodeIDLen)
	binary.Read(r, binary.LittleEndian, &cmd.NodeID)

	return cmd
}

func parseLeaveCommand(r io.Reader) *CommandLeave {
	cmd := &CommandLeave{}
	var nodeIdLen int32
	binary.Read(r, binary.LittleEndian, &nodeIdLen)
	cmd.NodeID = make([]byte, nodeIdLen)
	binary.Read(r, binary.LittleEndian, &cmd.NodeID)
	return cmd
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)
	return cmd
}

func parseDelCommand(r io.Reader) *CommandDel {
	cmd := &CommandDel{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)
	return cmd
}

func parseHasCommand(r io.Reader) *CommandHas {
	cmd := &CommandHas{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)
	return cmd
}
