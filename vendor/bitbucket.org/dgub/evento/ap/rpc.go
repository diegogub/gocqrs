package ap

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func InitRPC(port string) error {
	events := new(Streamer)
	rpc.Register(events)
	rpc.HandleHTTP()
	log.Println(port)
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		panic(e)
	}
	go http.Serve(l, nil)
	return nil
}
