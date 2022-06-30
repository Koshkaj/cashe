package main

import (
	"flag"
	"log"

	"github.com/Koshkaj/cashe/cache"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		LeaderAddr: *leaderAddr,
		IsLeader:   len(*leaderAddr) == 0,
	}

	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())

}
