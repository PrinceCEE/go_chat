package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type userHandler struct {
	services services.Services
}

func (h *userHandler) getAccount(c *gin.Context) error {
	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    "unauthorized user",
			Err:        errors.New("unauthorized user"),
			StatusCode: http.StatusUnauthorized,
		}
	}
	user := val.(*models.User)
	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "fetched user successfully",
		Data:    map[string]any{"user": user},
	})
	return nil
}

type UpdateAccountDto struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,alpha"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,alpha"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

func (h *userHandler) updateAccount(c *gin.Context) error {
	userId, ok := c.Params.Get("id")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	if user.ID != userId {
		return &utils.ServerError{
			Message:    utils.ErrUnauthorized.Error(),
			Err:        utils.ErrUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}
	}

	var updateDto UpdateAccountDto
	err := c.BindJSON(&updateDto)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if updateDto.Email != nil {
		user.Email = strings.ToLower(*updateDto.Email)
	}
	if updateDto.FirstName != nil {
		user.FirstName = *updateDto.FirstName
	}
	if updateDto.LastName != nil {
		user.LastName = *updateDto.LastName
	}

	userService := h.services.GetUserService()
	err = userService.UpdateUser(user, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Message: "user updated successfully",
		Success: true,
		Data:    map[string]any{"user": user},
	})
	return nil
}

func (h *userHandler) deleteAccount(c *gin.Context) error {
	userId, ok := c.Params.Get("id")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	if user.ID != userId {
		return &utils.ServerError{
			Message:    utils.ErrUnauthorized.Error(),
			Err:        utils.ErrUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}
	}

	userService := h.services.GetUserService()
	err := userService.DeleteUser(userId, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Message: "user deleted successfully",
		Success: true,
	})
	return nil
}
