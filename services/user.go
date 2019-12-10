package services

import (
	"fmt"
	"user/models"
	"user/repositories"
	"user/utils"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UserServices interface {
	CreateUser(*models.Body) (models.Token, error)
}

type user struct {
	repositories.DB
}

func New(repo repositories.DB) UserServices {

	return &user{
		repo,
	}
}

func (u *user) CreateUser(body *models.Body) (token models.Token, err error) {
	validateBody := &struct {
		UserName string `validate:"required"`
		Password string `validate:"required"`
		email    string `validate:"required,email"`
	}{
		body.UserName,
		body.Password,
		body.Email,
	}
	validate := validator.New()

	if errors := validate.Struct(validateBody); errors != nil {
		errStr := ""

		for idx, errs := range errors.(validator.ValidationErrors) {
			if idx == 0 {
				errStr = fmt.Sprintf("bad request: %s, ", errs.Namespace())
			} else {
				errStr = fmt.Sprintf("%s%s ,", errStr, errs.Namespace())
			}
		}
	}

	userID := (uuid.New()).String()
	password, err := utils.Hash(body.Password)
	err = u.InsertUser(&models.User{
		ID:       userID,
		Email:    body.Email,
		Password: password,
		UserName: body.UserName,
	})

	if err != nil {
		token = models.Token("")

		return
	}

	token, err = utils.Token(userID, "all")

	return
}