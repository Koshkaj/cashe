package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/hashicorp/raft"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
		raftAddr   = flag.String("raftaddr", ":4000", "listen address of raft server")
		nodeID     = flag.String("id", "", "node id of cashe/raft server")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr:       *listenAddr,
		LeaderAddr:       *leaderAddr,
		IsLeader:         len(*leaderAddr) == 0,
		RaftAddr:         *raftAddr,
		NodeID:           *nodeID,
		EvictionInterval: 5 * time.Second,
	}

	server := NewServer(opts, cache.New())
	go func() {
		server.logger.Fatal(server.Start())
	}()
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, os.Kill)
	<-terminate
	func() {
		_, leaderID := server.raft.LeaderWithID()
		if server.NodeID == string(leaderID) && len(server.members) > 0 {
			// If leader is disconnected and it has multiple members
			if err := server.raft.RemoveServer(raft.ServerID(server.NodeID), 0, 0).Error(); err != nil {
				server.logger.Errorf("Failed to remove old leader from Raft configuration: %s", err)
			}
		} else {
			// if member is disconnected, update config so that leader does not try to connect to it infinitely
			// its member, we can be sure that leaderaddr is provided
			conn, err := net.Dial("tcp", server.LeaderAddr)
			if err != nil {
				server.logger.Fatal(err)
			}
			server.writeLeaveCmd(conn)

		}
	}()
	server.logger.Infof("node shutdown [%s]", server.NodeID)
}
