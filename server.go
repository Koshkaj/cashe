package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/Koshkaj/cashe/client"
	"github.com/Koshkaj/cashe/core"
	rf "github.com/Koshkaj/cashe/raft"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

const raftTimeout = 5

type ServerOpts struct {
	ListenAddr       string
	IsLeader         bool
	LeaderAddr       string
	RaftAddr         string
	NodeID           string
	EvictionInterval time.Duration
}

type Server struct {
	ServerOpts
	members map[*client.Client]struct{}
	cache   cache.Cacher
	logger  *zap.SugaredLogger
	raft    *rf.RaftServer
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	l, _ := zap.NewProduction()
	lsugar := l.Sugar()
	fsm := cache.NewCacheFSM(c)
	r := rf.New(opts.NodeID, opts.RaftAddr, fsm)
	return &Server{
		ServerOpts: opts,
		cache:      c,
		members:    make(map[*client.Client]struct{}),
		logger:     lsugar,
		raft:       r,
	}
}

func (s *Server) EvictionLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	c := s.cache.(*cache.Cache)
	go func() {
		for range ticker.C {
			for key, value := range c.Expiry {
				if time.Now().After(value) {
					c.Delete([]byte(key))
					s.raft.Snapshot()
				}
			}
		}
	}()
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

	s.logger.Infow("server starting on port", "port", s.ListenAddr)
	s.EvictionLoop(s.ServerOpts.EvictionInterval)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.readLoop(conn)
	}
}

func (s *Server) writeJoinCmd(conn net.Conn) error {
	cmd := &core.CommandJoin{
		RaftAddr: []byte(s.RaftAddr),
		NodeID:   []byte(s.NodeID),
	}
	_, err := conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) writeLeaveCmd(conn net.Conn) error {
	cmd := &core.CommandLeave{
		NodeID: []byte(s.NodeID),
	}
	_, err := conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}
	s.logger.Infow("connected to leader ", "port", s.LeaderAddr)
	s.writeJoinCmd(conn)
	s.readLoop(conn)
	return nil
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
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
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *core.CommandSet:
		s.handleSetCommand(conn, v)
	case *core.CommandGet:
		s.handleGetCommand(conn, v)
	case *core.CommandJoin:
		s.handleJoinCommand(conn, v)
	case *core.CommandLeave:
		s.handleLeaveCommand(conn, v)
	case *core.CommandDel:
		s.handleDelCommand(conn, v)
	case *core.CommandHas:
		s.handleHasCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *core.CommandJoin) error {
	s.logger.Infow("member joined cluster", "address", conn.RemoteAddr())
	// raft server connect
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.logger.Errorf("failed to get raft configuration: %v", err)
		return err
	}
	var (
		nodeID = string(cmd.NodeID)
		addr   = string(cmd.RaftAddr)
	)
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.logger.Infof("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}
			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}
	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	s.members[client.NewFromConn(conn)] = struct{}{}
	s.logger.Debugf("node %s at %s joined successfully", nodeID, addr)
	return nil

}

func (s *Server) handleLeaveCommand(conn net.Conn, cmd *core.CommandLeave) error {
	s.logger.Infow("member leaving cluster", "address", conn.RemoteAddr())
	// raft server disconnect
	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.logger.Errorf("failed to get raft configuration: %v", err)
		return err
	}
	nodeID := string(cmd.NodeID)
	future := s.raft.RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := future.Error(); err != nil {
		return fmt.Errorf("error removing existing node %s", nodeID)
	}
	s.logger.Debugf("node %s at %s disconnected successfully", nodeID)
	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *core.CommandSet) error {
	s.logger.Infof("SET : %s => %s\n", cmd.Key, cmd.Value)
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	resp := core.ResponseSet{}
	f := s.raft.Apply(cmd.Bytes(), raftTimeout)
	if err := f.Error(); err != nil {
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

func (s *Server) handleDelCommand(conn net.Conn, cmd *core.CommandDel) error {
	s.logger.Infof("DEL %s", cmd.Key)

	resp := core.ResponseDel{}
	f := s.raft.Apply(cmd.Bytes(), raftTimeout)
	if err := f.Error(); err != nil {
		resp.Status = core.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = core.StatusOK

	_, err := conn.Write(resp.Bytes())
	return err
}

func (s *Server) handleHasCommand(conn net.Conn, cmd *core.CommandHas) error {
	resp := core.ResponseHas{}
	value := s.cache.Has(cmd.Key)
	if !value {
		resp.Status = core.StatusKeyNotFound
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Value = value
	resp.Status = core.StatusOK
	_, err := conn.Write(resp.Bytes())
	return err

}
