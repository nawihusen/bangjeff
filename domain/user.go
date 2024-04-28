package domain

import (
	"context"
	"time"
)

type User struct {
	ID       int64  `json:"id" form:"id" `
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Name     string `json:"name" form:"name"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
	Address  string `json:"address" form:"address"`
	DtmCrt   time.Time
}

type UserResponse struct {
	ID       int64     `json:"id" form:"id" `
	Username string    `json:"username" form:"username"`
	Name     string    `json:"name" form:"name"`
	Phone    string    `json:"phone" form:"phone"`
	Email    string    `json:"email" form:"email"`
	Address  string    `json:"address" form:"address"`
	DtmCrt   time.Time `json:"dtm_crt"`
}

type Options struct {
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
	Sort     string `json:"sort"`
	Order    string `json:"order"`
	Name     string `json:"name" `
	Username string `json:"username"`
}

type Metadata struct {
	TotalData int64  `json:"total_data"`
	TotalPage int64  `json:"total_page"`
	Page      int64  `json:"page"`
	Limit     int64  `json:"limit"`
	Sort      string `json:"sort"`
	Order     string `json:"order"`
}

type Token struct {
	Token  string `json:"token"`
	Active string `json:"active"`
}

type Test struct {
	Test string `json:"test,omitempty" form:"test,omitempty" `
}

type UsersResponse struct {
	Data     []UserResponse `json:"data"`
	Metadata Metadata       `json:"metadata"`
}

type UserUsecase interface {
	SignIn(ctx context.Context, user User) (token Token, err error)
	SignUp(ctx context.Context, user User) (err error)
	GetUsers(ctx context.Context, opt Options, meta *Metadata) (users []User, err error)
	GetUser(ctx context.Context, id int64) (users User, err error)
	GetByToken(ctx context.Context, token string) (users User, err error)
}

type UserRepository interface {
	SignUp(ctx context.Context, user User) (err error)
	GetUsers(ctx context.Context, opt Options) (users []User, err error)
	CountUsers(ctx context.Context, opt Options) (sum int64, err error)
	GetUser(ctx context.Context, id int64) (users User, err error)
	GetByUsername(ctx context.Context, username string) (user User, err error)
	SaveToken(ctx context.Context, id int64, token Token, active int64) (err error)
	GetByToken(ctx context.Context, token string) (users User, active int64, err error)
}
