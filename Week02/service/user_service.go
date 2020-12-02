package service

import (
	"Go-000/Week02/Dao"
	"Go-000/Week02/Data"
)

type UserService struct{
	d *Dao.Dao
}

func NewUserService()*UserService{
	return &UserService{
		d: Dao.NewDao(),
	}
}

func(u *UserService)GetUser(id int)(*Data.ServiceUser, error){
	d, err := u.d.GetUser(id)
	if err != nil{
		return nil, err
	}

	return &Data.ServiceUser{
		Name: d.Name,
		Age:  d.Age,
	}, nil
}
