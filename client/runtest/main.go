package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Koshkaj/cashe/client"
)

func main() {
	SendCommand()
}

func SendCommand() {
	c, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		var (
			key   = []byte(fmt.Sprintf("key_%d", i))
			value = []byte(fmt.Sprintf("value_%d", i))
		)
		if err := c.Set(context.Background(), key, value, 2); err != nil {
			log.Fatal(err)
		}
		val, _ := c.Get(context.Background(), key)
		fmt.Println(val)
		// ok, _ := c.Has(context.Background(), key)
		// fmt.Println(ok)

		time.Sleep(time.Second)
		// err = c.Delete(context.Background(), key)
	}
	c.Close()
}
