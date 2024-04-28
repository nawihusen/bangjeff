package handler

import (
	"bangjeff/domain"
	"bangjeff/helper"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// "github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

// Handler is REST API handler for Service System
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func (uh *UserHandler) SignUp(c *fiber.Ctx) error {
	var input domain.User
	err := c.BodyParser(&input)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Error(err)
		return c.Status(400).SendString("username and password cant blank")
	}

	err = uh.UserUsecase.SignUp(c.Context(), input)
	if err != nil {
		log.Error(err)
		return c.Status(fasthttp.StatusBadRequest).SendString(err.Error())
		// return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return c.Status(fasthttp.StatusOK).JSON("Pendaftaran Berhasil")
}

func (uh *UserHandler) SignIn(c *fiber.Ctx) error {
	var input domain.User
	err := c.BodyParser(&input)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		log.Error(err)
		return c.Status(400).SendString("username and password cant blank")
	}

	token, err := uh.UserUsecase.SignIn(c.Context(), input)
	if err != nil {
		log.Error(err)
		return c.Status(fasthttp.StatusBadRequest).SendString(err.Error())
		// return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	return c.Status(fasthttp.StatusOK).JSON(token)
}

func (uh *UserHandler) GetUsers(c *fiber.Ctx) error {
	options := domain.Options{}
	meta := domain.Metadata{}
	page, err := strconv.ParseInt(c.Query("page", fmt.Sprintf("%d", viper.GetInt("default_page"))), 10, 64)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	options.Page = page
	meta.Page = page

	limit, err := strconv.ParseInt(c.Query("limit", fmt.Sprintf("%d", viper.GetInt("default_limit"))), 10, 64)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}
	options.Limit = limit
	meta.Limit = limit

	sort := c.Query("sort", "id")
	options.Sort = sort
	meta.Sort = sort

	order := c.Query("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	options.Order = order
	meta.Order = order

	name := c.Query("name")
	if name != "" {
		options.Name = name
	}

	username := c.Query("username")
	if username != "" {
		options.Username = username
	}

	users, err := uh.UserUsecase.GetUsers(c.Context(), options, &meta)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusInternalServerError)
	}

	data := toUserResponseBulk(users)

	response := domain.UsersResponse{
		Data:     data,
		Metadata: meta,
	}

	return c.Status(fasthttp.StatusOK).JSON(response)
}

func (uh *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		log.Error(err)
		return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	user, err := uh.UserUsecase.GetUser(c.Context(), id)
	if err != nil {
		log.Error(err)
		return c.Status(400).SendString(err.Error())
		// return helper.HTTPSimpleResponse(c, fasthttp.StatusBadRequest)
	}

	response := toUserResponse(user)

	return c.Status(fasthttp.StatusOK).JSON(response)
}

func (uh *UserHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(domain.User)

	resp := toUserResponse(user)

	return c.Status(fasthttp.StatusOK).JSON(resp)
}

func (uh *UserHandler) Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		auth := c.GetReqHeaders()
		if len(auth["Authorization"]) < 1 {
			return c.Status(401).SendString("unauthorized")
		}
		authorization := auth["Authorization"][0]
		if authorization != "" {
			token = authorization[7:]
		}
		if token == "" {
			log.Warning("Invalid token")
			return helper.HTTPSimpleResponse(c, fasthttp.StatusUnauthorized)
		}

		user, err := uh.UserUsecase.GetByToken(c.Context(), token)
		if err != nil {
			log.Error(err)
			return c.Status(401).SendString(err.Error())
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func toUserResponse(user domain.User) domain.UserResponse {
	return domain.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Address:  user.Address,
		DtmCrt:   user.DtmCrt,
	}
}

func toUserResponseBulk(users []domain.User) []domain.UserResponse {
	response := []domain.UserResponse{}

	for _, v := range users {
		temp := toUserResponse(v)
		response = append(response, temp)
	}

	return response
}
