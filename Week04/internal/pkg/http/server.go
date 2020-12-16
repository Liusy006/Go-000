package http

import (
	"context"
	"net/http"
)

type GirlsInterface interface {
	RequestGirls(id string)string
}

type HttpServer struct{
	addr string
	sm *http.ServeMux
	server http.Server

	gi GirlsInterface
}

func NewHttpServer(addr string)*HttpServer{
	return &HttpServer{
		addr:addr,
		sm: http.NewServeMux(),
	}
}

func (s *HttpServer)Run()error{
	s.sm.HandleFunc("/girls", func(rsp http.ResponseWriter, req *http.Request){
		if s.gi != nil{
			s.gi.RequestGirls("")
		}
	})

	s.server = http.Server{
		Addr:              ":8080",
		Handler:           s.sm,
	}
	return s.server.ListenAndServe()
}

func (s *HttpServer)RegisterGirlsHandler( g GirlsInterface){
	s.gi = g
}

func (s *HttpServer)ShutDown()error{
	return s.server.Shutdown(context.Background())
}