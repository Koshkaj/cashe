package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Koshkaj/cashe/cache"
	"github.com/Koshkaj/cashe/client"
	"github.com/Koshkaj/cashe/core"
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
		for i := 0; i < 8; i++ {
			SendCommand()
			time.Sleep(time.Millisecond * 2)
		}
	}()

	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())

}

func SendCommand() {
	cmd := &core.CommandSet{
		Key:   []byte("keytest"),
		Value: []byte("valuetest"),
		TTL:   9,
	}
	cl, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cl.Set(context.Background(), []byte("foo"), []byte("bar"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(cmd.Bytes())
}
