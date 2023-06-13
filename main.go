package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Koshkaj/cashe/cache"
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
		log.Fatal(server.Start())
	}()
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, os.Kill)
	<-terminate
}
