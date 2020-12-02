package main

import (
	"Go-000/Week02/service"
	"database/sql"
	"errors"
	"fmt"
)

func main(){
	s := service.NewUserService()
	d, err := s.GetUser(1)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			fmt.Printf("%+v", err)
		}

		return
	}
	fmt.Println(d)
}
