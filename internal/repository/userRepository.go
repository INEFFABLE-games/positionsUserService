package repository

import (
	"context"
	"database/sql"
	"positionsUserService/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) Insert(ctx context.Context,user models.User)error{
	_,err := u.db.ExecContext(ctx,"insert into users(uid,login,password,balance) values($1,$2,$3,$4)",user.Uid,user.Login,user.Password,0)
	return err
}

func (u *UserRepository) UpdateBalance(ctx context.Context,uid string,balance int64)error{
	_,err := u.db.ExecContext(ctx,"update users set balance = $1 where uid = $2",balance,uid)
	return err
}

func (u *UserRepository) GetBalance(ctx context.Context,uid string)(int64,error){
	var balance int64
	err := u.db.QueryRowContext(ctx,"select balance from users where uid = $1",uid).Scan(&balance)
	return balance,err
}

func (u *UserRepository) GetUserByLogin(ctx context.Context,login string)(models.User,error){
	var user models.User
	err := u.db.QueryRowContext(ctx,"select uid,login,password,balance from users where login = $1",login).Scan(&user.Uid,&user.Login,&user.Password,&user.Balance)
	return user,err
}

func NewUserRepository(db *sql.DB)*UserRepository{
	return &UserRepository{db: db}
}
