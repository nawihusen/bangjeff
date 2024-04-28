package usecase

import (
	"bangjeff/domain"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	// log "github.com/sirupsen/logrus"
)

// userUsecase is struct usecase
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUsersecase is constructor of account usecase
func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) SignUp(ctx context.Context, user domain.User) (err error) {
	_, err = uu.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		if err.Error() == "username not found" {
			err = nil
		} else {
			return err
		}
	} else if err == nil {
		return errors.New("username already exist")
	}

	// Default
	if user.Name == "" {
		user.Name = "bangjeff"
	}
	if user.Phone == "" {
		user.Phone = "00000000"
	}
	if user.Email == "" {
		user.Email = "bangjeff@gmail.com"
	}
	if user.Address == "" {
		user.Address = "Situbondo"
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	password := string(hashPass)
	user.Password = password

	err = uu.userRepo.SignUp(ctx, user)

	return err
}

func (uu *userUsecase) SignIn(ctx context.Context, user domain.User) (token domain.Token, err error) {
	result, err := uu.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		return token, err
	}

	errPw := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if errPw != nil {
		return token, errors.New("wrong password")
	}

	date := time.Now().Format(time.RFC3339)
	tokenByte := sha256.Sum256([]byte(strconv.FormatInt(result.ID, 10) + "_" + result.Username + date))
	token.Token = base64.URLEncoding.EncodeToString(tokenByte[:])

	timeLimit := time.Now().Add(24 * time.Hour)

	active := timeLimit.Unix()

	uu.userRepo.SaveToken(ctx, result.ID, token, active)

	token.Active = fmt.Sprintf("%v", timeLimit)

	return token, err
}

func (uu *userUsecase) GetUsers(ctx context.Context, opt domain.Options, meta *domain.Metadata) (users []domain.User, err error) {
	total, err := uu.userRepo.CountUsers(ctx, opt)
	if err != nil {
		return
	}

	meta.TotalData = total

	users, err = uu.userRepo.GetUsers(ctx, opt)

	meta.TotalPage = countTotalPage(total, meta.Limit)

	return users, err
}

func (uu *userUsecase) GetUser(ctx context.Context, id int64) (user domain.User, err error) {
	user, err = uu.userRepo.GetUser(ctx, id)
	return user, err
}

func (uu *userUsecase) GetByToken(ctx context.Context, token string) (user domain.User, err error) {
	user, active, err := uu.userRepo.GetByToken(ctx, token)
	if err != nil {
		return user, err
	}

	expire := time.Now().Unix()

	if active <= expire {
		err = errors.New("token had expire")
	}

	return user, err
}

func countTotalPage(totalData, limit int64) int64 {
	totalpage := totalData / limit
	left := totalData % limit
	if left > 0 {
		totalpage += 1
	}
	return totalpage
}
