package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/Koshkaj/cashe/core"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	//	followers map[net.Conn]struct{}
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listening error: %s", err)
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

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection made: ", conn.RemoteAddr())
	for {
		cmd, err := core.ParseCommand(conn)
		if err != nil {
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *core.CommandSet:
		s.handleSetCommand(conn, v)
	case *core.CommandGet:
		s.handleGetCommand(conn, v)
		// case *CommandDel:
		// 	s.handleDelCommand(conn, v)
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *core.CommandSet) error {
	return s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)*time.Second)
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *core.CommandGet) ([]byte, error) {
	return s.cache.Get(cmd.Key)
}
