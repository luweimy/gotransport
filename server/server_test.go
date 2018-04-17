package server

import (
	"log"
	"testing"

	"github.com/luweimy/gotransport"
	"github.com/luweimy/gotransport/protocol"
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

func onErr(transport gotransport.Transport, err error) {
	log.Println("ON-ERR", err)
	transport.WriteString("!>err\n => ")
}

func onClose(transport gotransport.Transport) {
	log.Println("ON-CLOSE", transport.Peer())
}

func onMessage(transport gotransport.Transport, packet protocol.Protocol) {
	log.Println("ON-MSG", transport.Peer(), packet.Type(), string(packet.Payload()))
	transport.WriteString("=> ok")
}

func TestNew(t *testing.T) {
	server := New(gotransport.WithConnect(onConnect), gotransport.WithClose(onClose), gotransport.WithError(onErr), gotransport.WithMessage(onMessage), gotransport.WithFactory(protocol.LineFactory{}))
	err := server.Listen("tcp", "127.0.0.1:9090")
	errorCheck(err)
}
