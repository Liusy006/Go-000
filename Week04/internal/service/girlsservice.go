package service

import "Go-000/Week04/internal/data"

type GirlsDaoInterface interface {
	Save(girls *data.Girls)string
}
type GirlsService struct{
	gd GirlsDaoInterface
}

func NewGirlsService(i GirlsDaoInterface)*GirlsService{
	return &GirlsService{gd:i}
}

func (g *GirlsService)Save(girls *data.Girls){
	g.gd.Save(girls)
}
