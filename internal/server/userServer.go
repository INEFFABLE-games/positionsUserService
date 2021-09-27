package server

import (
	"context"
	authService "github.com/INEFFABLE-games/authService/models"
	"github.com/INEFFABLE-games/positionsUserService/internal/models"
	"github.com/INEFFABLE-games/positionsUserService/internal/service"
	"github.com/INEFFABLE-games/positionsUserService/protocol"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"strings"
)

type UserServer struct {
	userService *service.UserService

	protocol.UnimplementedUserServiceServer
}

func middlewareRTValidation(token string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(
		strings.TrimPrefix(token, "Bearer "),
		&authService.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("operationToken"), nil
		},
	)
	if err != nil {
		return "", err
	}

	claim, ok := parsedToken.Claims.(*authService.CustomClaims)
	if !ok {
		return "", err
	}

	return claim.Uid, err
}

func (u *UserServer) Create(ctx context.Context, request *protocol.CreateRequest) (*protocol.CreateReply, error) {

	user := models.User{
		Uid:      "",
		Login:    request.GetLogin(),
		Password: request.GetPassword(),
		Balance:  0,
	}

	err := u.userService.Create(ctx, user)
	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "userServer",
			"action ":  "create user",
		}).Errorf("unable to create user %v", err.Error())

		mes := "unable to create user"
		return &protocol.CreateReply{Message: &mes}, err
	}

	mes := "user created"
	return &protocol.CreateReply{Message: &mes}, err
}

func (u *UserServer) RefreshTokens(ctx context.Context, request *protocol.RefreshTokensRequest) (*protocol.RefreshTokensReply, error) {

	rt := request.GetRT()

	uid, err := middlewareRTValidation(rt)
	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "userServer",
			"action ":  "validate RT",
		}).Errorf("unable to validate RT %v", err.Error())

		return &protocol.RefreshTokensReply{
			RT:  nil,
			JWT: nil,
		}, err
	}

	jwtt, rt, err := u.userService.Refresh(ctx, uid)

	return &protocol.RefreshTokensReply{
		RT:  &rt,
		JWT: &jwtt,
	}, err
}

func (u *UserServer) UpdateBalance(ctx context.Context, request *protocol.UpdateBalanceRequest) (*protocol.UpdateBalanceReply, error) {

	uid := request.GetUid()
	balance := request.GetBalance()

	err := u.userService.UpdateBalance(ctx, uid, balance)
	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "userServer",
			"action ":  "update balance",
		}).Errorf("unable to update balance %v", err.Error())
	}

	return &protocol.UpdateBalanceReply{}, err
}
func (u *UserServer) GetBalance(ctx context.Context, request *protocol.GetBalanceRequest) (*protocol.GetBalanceReply, error) {
	uid := request.GetUid()

	balance, err := u.userService.GetBalance(ctx, uid)
	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "userServer",
			"action ":  "get balance",
		}).Errorf("unable to get balance %v", err.Error())
	}

	return &protocol.GetBalanceReply{Balance: &balance}, err
}

func (u *UserServer) Login(ctx context.Context, request *protocol.LoginRequest) (*protocol.LoginReply, error) {

	login := request.GetLogin()
	password := request.GetPassword()

	JWT, RT, err := u.userService.Login(ctx, login, password)
	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "userServer",
			"action ":  "login user",
		}).Errorf("unable to login %v", err.Error())
	}

	return &protocol.LoginReply{Jwt: &JWT, RT: &RT}, err
}

func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}
