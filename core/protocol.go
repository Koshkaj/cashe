package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// type Command string

// const (
// 	CMDSet Command = "SET"
// 	CMDGet Command = "GET"
// )

// type Message struct {
// 	Cmd   Command
// 	Key   []byte
// 	Value []byte
// 	TTL   time.Duration
// }

// func (m *Message) ToBytes() []byte {
// 	switch m.Cmd {
// 	case CMDSet:
// 		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL)
// 		return []byte(cmd)
// 	case CMDGet:
// 		cmd := fmt.Sprintf("%s %s ", m.Cmd, m.Key)
// 		return []byte(cmd)
// 	default:
// 		panic("unknown command")
// 	}
// }

// func parseMessage(raw []byte) (*Message, error) {
// 	var (
// 		rawStr = string(raw)
// 		parts  = strings.Split(rawStr, " ")
// 	)
// 	if len(parts) < 2 {
// 		return nil, fmt.Errorf("invalid command %s", raw)
// 	}
// 	msg := &Message{
// 		Cmd: Command(parts[0]),
// 		Key: []byte(parts[1]),
// 	}
// 	if msg.Cmd == CMDSet {
// 		if len(parts) < 4 {
// 			return nil, fmt.Errorf("invalid SET command")
// 		}
// 		msg.Value = []byte(parts[2])
// 		ttl, err := strconv.Atoi(parts[3])
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid SET TTL")
// 		}
// 		msg.TTL = time.Duration(ttl) * time.Minute
// 	}
// 	return msg, nil
// }

type Command byte

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
)

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

type CommandDel struct {
	Key []byte
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	// bufreader := bufio.NewReader(r)
	binary.Read(r, binary.LittleEndian, &cmd)
	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	default:
		return nil, fmt.Errorf("invalid command %s", string(cmd))
	}
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
