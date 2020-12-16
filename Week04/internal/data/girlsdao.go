package data

type Girls struct{
	Id string
	Visit_times int
}

type GirlsDao struct{

}

func NewGirlsDao()*GirlsDao{
	return &GirlsDao{
	}
}

func (g *GirlsDao)Save(gs *Girls)string{
	//真正的业务操作
	return ""
}
