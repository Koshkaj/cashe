package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/Koshkaj/cashe/client"
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

	go func() {
		time.Sleep(time.Second * 2)
		cl, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 8; i++ {
			SendCommand(cl)
		}
		cl.Close()
	}()

	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())

}

func SendCommand(c *client.Client) {
	_, err := c.Set(context.Background(), []byte("anyhow"), []byte("wassup"), 20)
	if err != nil {
		log.Fatal(err)
	}
}
