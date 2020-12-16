package biz

import (
	"Go-000/Week04/internal/data"
	"Go-000/Week04/internal/service"
)

type GirlsHandler struct{
	s *service.GirlsService
}

func NewGirlsHandler(s *service.GirlsService)*GirlsHandler{
	return &GirlsHandler{s:s}
}

func (g *GirlsHandler)RequestGirls(id string)string{
	g.s.Save(&data.Girls{
		Id:id,
	})

	return ""
}

