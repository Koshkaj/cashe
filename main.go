package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
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
		if opts.IsLeader {
			SendCommand()
		}
	}()

	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())

}

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	io.ReadFull(rand.Reader, buf)
	return buf
}

func SendCommand() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			c, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Second * 2)
			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("value_%d", i))
			)
			err = c.Set(context.Background(), key, value, 20)
			if err != nil {
				log.Fatal(err)
			}

			_, err = c.Get(context.Background(), key)
			c.Close()
		}(i)
	}
}
