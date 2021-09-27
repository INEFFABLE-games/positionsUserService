package service

import (
	"context"
	authProto "github.com/INEFFABLE-games/authService/protocol"
	"github.com/INEFFABLE-games/positionsUserService/internal/client"
	"github.com/INEFFABLE-games/positionsUserService/internal/models"
	"github.com/INEFFABLE-games/positionsUserService/internal/repository"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	authClient     *client.AuthClient
}

func (u *UserService) Create(ctx context.Context, user models.User) error {

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	user.Uid = uid.String()

	cryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(cryptPass)
	return u.userRepository.Insert(ctx, user)
}

func (u *UserService) UpdateBalance(ctx context.Context, uid string, balance int64) error {
	return u.userRepository.UpdateBalance(ctx, uid, balance)
}

func (u *UserService) GetBalance(ctx context.Context, uid string) (int64, error) {
	return u.userRepository.GetBalance(ctx, uid)
}

func (u *UserService) Refresh(ctx context.Context, uid string) (string, string, error) {
	reply, err := u.authClient.Refresh(ctx, &authProto.RefreshRequest{
		Uid: &uid,
	})
	if err != nil {
		return "", "", err
	}

	jwt := reply.GetJwt()
	Rt := reply.GetRt()

	return jwt, Rt, err
}

func (u *UserService) Login(ctx context.Context, login string, password string) (string, string, error) {

	user, err := u.userRepository.GetUserByLogin(ctx, login)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", err
	}

	reply, err := u.authClient.Refresh(ctx, &authProto.RefreshRequest{
		Uid: &user.Uid,
	})
	if err != nil {
		return "", "", err
	}

	jwt := reply.GetJwt()
	Rt := reply.GetRt()

	return jwt, Rt, err
}

func NewUserService(userRepository *repository.UserRepository, authClient *client.AuthClient) *UserService {
	return &UserService{userRepository: userRepository, authClient: authClient}
}
