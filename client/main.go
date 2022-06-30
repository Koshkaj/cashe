package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("SET foo bar 25"))
	time.Sleep(time.Second * 2)
	conn.Write([]byte("GET foo"))
	buf := make([]byte, 1000)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))

	select {}
}
