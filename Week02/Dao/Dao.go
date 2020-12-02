package Dao

import (
	"Go-000/Week02/Data"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

type Dao struct{

}

func NewDao()*Dao{
	return &Dao{}
}

func (d *Dao)GetUser(id int)( *Data.DaoUser, error){
	return nil, errors.Wrap(sql.ErrNoRows, fmt.Sprintf("user:%d is not exist", id))
}