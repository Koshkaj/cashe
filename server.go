package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/Koshkaj/cashe/client"
	"github.com/Koshkaj/cashe/core"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	members map[*client.Client]struct{}
	cache   cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		members:    make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listening error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.readLoop(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}
	log.Println("connected to leader : ", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, core.CmdJoin)
	s.readLoop(conn)
	return nil
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection made: ", conn.RemoteAddr())
	for {
		cmd, err := core.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
	// fmt.Println("connection closed: ", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *core.CommandSet:
		s.handleSetCommand(conn, v)
	case *core.CommandGet:
		s.handleGetCommand(conn, v)
	case *core.CommandJoin:
		s.handleJoinCommand(conn, v)

		// case *CommandDel:
		// 	s.handleDelCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *core.CommandJoin) error {
	fmt.Println("member joined cluster", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}
	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *core.CommandSet) error {
	log.Printf("SET : %s => %s\n", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error", err)
			}
		}
	}()

	resp := core.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)*time.Second); err != nil {
		resp.Status = core.StatusError
		_, err := conn.Write(resp.Bytes())
		return err

	}
	resp.Status = core.StatusOK

	_, err := conn.Write(resp.Bytes())
	return err
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *core.CommandGet) error {
	resp := core.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = core.StatusKeyNotFound
		_, err := conn.Write(resp.Bytes())
		return err

	}
	resp.Value = value
	resp.Status = core.StatusOK
	_, err = conn.Write(resp.Bytes())
	return err
}
