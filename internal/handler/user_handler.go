package handler

import (
	"errors"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/response"
	"github.com/H0wZy/authsys/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func (ctrl *UserHandler) Create(ctx *gin.Context) {
	var input dto.CreateUser

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, "error while binding json:", err.Error())
		return
	}

	user, err := ctrl.userService.Create(ctx.Request.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists), errors.Is(err, service.ErrUsernameAlreadyExists):
			response.Conflict(ctx, err.Error())

		case errors.Is(err, service.ErrUsernameCantContainAt), errors.Is(err, service.ErrInvalidBirthDate):
			response.BadRequest(ctx, err.Error())

		default:
			response.InternalServerError(ctx)
		}

		return
	}

	response.Created(ctx, dto.ToUserResponse(user))
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
