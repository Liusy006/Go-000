package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	l         *net.TCPListener
	closed    bool
	closeChan chan struct{}

	connLastRWTime map[*net.TCPConn]time.Time
	mu             sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		closed:         false,
		closeChan:      make(chan struct{}),
		connLastRWTime: make(map[*net.TCPConn]time.Time),
	}
}

func (s *Server) Serve(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	go s.checkIdleConn()

	s.l = l

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Millisecond)
				continue
			}
			return err
		}

		s.connLastRWTime[conn] = time.Now()

		msgChan := make(chan []byte, 10)
		go s.read(conn, msgChan)
		go s.write(conn, msgChan)
	}
}

func (s *Server) Close() {
	if s.closed {
		return
	}

	s.l.Close()
	s.closed = true
	s.closeAllConn()
	s.closeChan <- struct{}{}
}

func (s *Server) read(conn *net.TCPConn, msgChan chan []byte) {
	defer s.removeConn(conn)
	reader := bufio.NewReader(conn)
	for {
		var bodyLen int32
		err := binary.Read(reader, binary.BigEndian, &bodyLen)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read failed, err: %v\n", err)
			}

			return
		}

		if bodyLen <= 0 {
			continue
		}
		s.connLastRWTime[conn] = time.Now()

		message := make([]byte, bodyLen)
		_, err = io.ReadFull(reader, message)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read failed, err: %v\n", err)
			}
			return
		}

		msgChan <- message
	}
}

func (s *Server) write(conn *net.TCPConn, msgChan chan []byte) {
	for {
		msg := <-msgChan
		if len(msg) <= 0 {
			return
		}

		//handle message
		func(msg []byte) []byte {
			fmt.Println(string(msg))
			return msg
		}(msg)

		bodyLen := len(msg)
		err := binary.Write(conn, binary.BigEndian, bodyLen)
		if err != nil {
			return
		}

		_, err = conn.Write(msg)
		if err != nil {
			return
		}
	}
}

func (s *Server) removeConn(conn *net.TCPConn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.connLastRWTime, conn)
}

func (s *Server) checkIdleConn() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if s.closed {
				return
			}
			s.closeIdleConn()
		case <-s.closeChan:
			return
		}

	}
}

func (s *Server) closeIdleConn() {
	curTime := time.Now()
	conns := make([]*net.TCPConn, 0)
	s.mu.RLock()
	for conn, readTime := range s.connLastRWTime {
		if curTime.Sub(readTime) > time.Second*10 {
			conns = append(conns, conn)
		}
	}
	s.mu.RUnlock()

	if len(conns) < 1 {
		return
	}

	s.mu.Lock()
	for _, conn := range conns {
		if readTime, ok := s.connLastRWTime[conn]; ok && curTime.Sub(readTime) > time.Second*10 {
			conn.Close()
			delete(s.connLastRWTime, conn)
		}
	}
	s.mu.Unlock()
}

func (s *Server) closeAllConn() {
	s.mu.Lock()
	for conn, _ := range s.connLastRWTime {
		conn.Close()
		delete(s.connLastRWTime, conn)
	}
	s.mu.Unlock()
}
