package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/repositories"
	"github.com/grigoriy-st/YL-Golang/pkg/response"
	"github.com/grigoriy-st/YL-Golang/pkg/websocket"
	"github.com/sirupsen/logrus"
)

type Agent struct {
	Route *gin.RouterGroup
}

func (a *Agent) Index(ctx *gin.Context) {
	agents, err := repositories.AgentRepository().GetAllAgents()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, agents)
}

func (a *Agent) WebSocket(ctx *gin.Context) {
	err := websocket.Connect(ctx)
	if err != nil {
		logrus.Errorln("Wrong websocket connection: %s", err.Error())
		response.BadRequest(ctx, "incorrect websocket")
		return
	}
}
