package client

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/luweimy/gotransport"
)

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func TestClient(t *testing.T) {
	client := New()
	client.Options(gotransport.WithFactory(gotransport.LineProtocol))
	client.Options(gotransport.WithMessage(func(transport gotransport.Transport, packet gotransport.Protocol) {
		fmt.Println(string(packet.Payload()))
	}))
	client.Options(gotransport.WithClosed(func(transport gotransport.Transport) {
		fmt.Println("on-close")
	}))
	for {
		time.Sleep(time.Second)
		err := client.Connect("tcp", "127.0.0.1:9090")
		if err != nil {
			log.Println("ERROR: conn err", err)
			continue
		}
		for {
			time.Sleep(time.Second * 2)
			//p := transport.NewlineProtocol()
			//p.Type = 0x01
			//p.Value = bytes.Repeat([]byte(fmt.Sprintf("hello,world %d", rand.Intn(10))), 1)
			//p.SetPayload([]byte("hello,world"))
			//n, err := p.WriteTo(client.conn)

			n, err := client.WriteString("hello,world")
			if err != nil {
				log.Println("ERROR: write to error", err)
				break
			}
			log.Println("INFO: write bytes", n, n/1024, "KB")
		}

	}
}
