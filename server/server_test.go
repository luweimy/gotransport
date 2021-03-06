package server

import (
	"log"
	"testing"

	"github.com/luweimy/gotransport"
)

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func onConnect(transport gotransport.Transport) bool {
	log.Println("ON-CONN", transport.Peer())
	transport.WriteString("=> ")
	return true
}

func onClosing(transport gotransport.Transport, err error) {
	log.Println("ON-CLOSE", transport.Peer(), err)
}

func onMessage(transport gotransport.Transport, packet gotransport.Protocol) {
	log.Println("ON-MSG", transport.Peer(), packet.FlagOptions(), string(packet.Payload()))
	//transport.WriteString(string(packet.Payload()))
	transport.WritePacket(packet)
}

func TestServer(t *testing.T) {
	server := New(gotransport.WithProtocol(gotransport.LineProtocol))
	server.Options(gotransport.WithConnected(onConnect))
	server.Options(gotransport.WithClosing(onClosing))
	server.Options(gotransport.WithMessage(onMessage))

	//go func() {
	//	time.Sleep(time.Second * 3)
	//	fmt.Println("server close before")
	//	server.Close()
	//	fmt.Println("server close after")
	//}()
	err := server.Listen("tcp", "127.0.0.1:9090")
	errorCheck(err)
}
