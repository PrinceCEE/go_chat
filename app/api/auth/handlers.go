package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type authHandler struct {
	services services.Services
}

type SignupDto struct {
	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName  string `json:"last_name" validate:"required,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (h *authHandler) signUp(c *gin.Context) error {
	var signupDto SignupDto
	err := c.BindJSON(&signupDto)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	userService := h.services.GetUserService()
	authService := h.services.GetAuthService()

	signupDto.Email = strings.ToLower(signupDto.Email)

	userExists, err := userService.GetUser(repositories.GetUserParams{
		Email: signupDto.Email,
	}, nil)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return &utils.ServerError{
				Message:    err.Error(),
				Err:        err,
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	if userExists != nil {
		return &utils.ServerError{
			Message:    "account already exists",
			Err:        utils.ErrDuplicateRecord,
			StatusCode: http.StatusBadRequest,
		}
	}

	tx, _ := h.services.GetDB().Begin(context.Background())

	user := &models.User{
		FirstName: signupDto.FirstName,
		LastName:  signupDto.LastName,
		Email:     signupDto.Email,
	}
	err = userService.CreateUser(user, tx)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	password, err := utils.GeneratePasswordHash(signupDto.Password)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    "internal server error",
			StatusCode: http.StatusInternalServerError,
		}
	}

	auth := &models.Auth{
		UserID:   user.ID,
		Password: password,
	}
	err = authService.CreateAuth(auth, tx)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return &utils.ServerError{
			Message:    "internal server error",
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	accessToken, _ := utils.GenerateToken(&utils.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
	})

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "sign up successful",
		Data: map[string]any{
			"user": user,
		},
		Meta: &utils.ResponseMeta{
			AccessToken: accessToken,
		},
	})

	return nil
}

type SigninDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *authHandler) signIn(c *gin.Context) error {
	var signinDto SigninDto
	err := c.BindJSON(&signinDto)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	userService := h.services.GetUserService()
	authService := h.services.GetAuthService()

	signinDto.Email = strings.ToLower(signinDto.Email)
	user, err := userService.GetUser(repositories.GetUserParams{
		Email: signinDto.Email,
	}, nil)
	if err != nil {
		se := &utils.ServerError{
			Err: err,
		}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			se.Message = "user not found"
			se.StatusCode = http.StatusNotFound
		default:
			se.Message = err.Error()
			se.StatusCode = http.StatusInternalServerError
		}

		return se
	}

	auth, err := authService.GetUserAuth(user.ID, nil)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = utils.ComparePassword(signinDto.Password, auth.Password)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    "invalid credential",
			StatusCode: http.StatusUnauthorized,
		}
	}

	accessToken, _ := utils.GenerateToken(&utils.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
	})
	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "sign in successful",
		Data:    map[string]any{"user": user},
		Meta:    &utils.ResponseMeta{AccessToken: accessToken},
	})

	return nil
}
