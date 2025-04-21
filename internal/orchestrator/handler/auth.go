package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/services"
	"github.com/grigoriy-st/YL-Golang/pkg/jwt"
	"github.com/grigoriy-st/YL-Golang/pkg/response"
)

type Auth struct {
	Route *gin.RouterGroup
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *Auth) Register(ctx *gin.Context) {
	var request AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "invalid data")
		return
	}

	_, err := services.UserService().Create(request.Login, request.Password)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Data(ctx, "successful")
}

func (a *Auth) Login(ctx *gin.Context) {
	var request AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "invalid data")
		return
	}

	userId := services.UserService().Authorization(request.Login, request.Password)

	if userId == 0 {
		response.BadRequest(ctx, "authorization failed")
		return
	}

	token, err := jwt.New().CreateUserToken(userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, map[string]interface{}{
		"token":   token.Token,
		"user_id": userId,
	})
}
