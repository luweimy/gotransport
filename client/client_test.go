package client

import (
	"log"
	"testing"
	"time"

	transport "github.com/luweimy/gotransport/protocol"
)

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	client := New()
	for {
		time.Sleep(time.Second)
		err := client.Connect("tcp", "127.0.0.1:9090")
		if err != nil {
			log.Println("ERROR: conn err", err)
			continue
		}
		for {
			time.Sleep(time.Second * 2)
			p := transport.NewLine()
			//p.Type = 0x01
			//p.Value = bytes.Repeat([]byte(fmt.Sprintf("hello,world %d", rand.Intn(10))), 1)
			p.SetPayload([]byte("hello,world"))
			n, err := p.WriteTo(client.conn)
			if err != nil {
				log.Println("ERROR: write to error", err)
				break
			}
			log.Println("INFO: write bytes", n, n/1024, "KB")
		}

	}
}
