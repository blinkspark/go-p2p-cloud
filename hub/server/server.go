package server

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Server struct {
	port     int
	portAddr string
	net.Listener
	connChan chan net.Conn
}

func NewServer(port int) (s *Server, err error) {
	s = &Server{
		port:     port,
		portAddr: fmt.Sprintf(":%d", port),
		connChan: make(chan net.Conn, 1),
	}
	s.Listener, err = net.Listen("tcp", s.portAddr)
	if err != nil {
		return nil, err
	}
	return
}

func (s *Server) Serv(ctx context.Context) {
	go func() {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				log.Println("listen err:", err)
				continue
			}
			s.connChan <- conn
		}

	}()
	for {
		select {
		case <-ctx.Done():
			s.Listener.Close()
			close(s.connChan)
			return
		default:
			continue
		}
	}
}
