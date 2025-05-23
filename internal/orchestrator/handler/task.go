package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/repositories"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/services"
	"github.com/grigoriy-st/YL-Golang/pkg/jwt"
	"github.com/grigoriy-st/YL-Golang/pkg/response"
)

type Task struct {
	Route *gin.RouterGroup
}

type TaskCreateRequest struct {
	Expression string `json:"expression" binding:"required"`
}

func (p *Task) Index(ctx *gin.Context) {
	userId := jwt.New().JwtUserId(ctx)
	tasks, err := repositories.TaskRepository().GetAllTasks(userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, tasks)
}

func (p *Task) Store(ctx *gin.Context) {
	userId := jwt.New().JwtUserId(ctx)
	var request TaskCreateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "невалидный expression")
		return
	}

	task, err := services.TaskService().Create(request.Expression, userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, task)
}

// @Summary Get task by id
// @Tags Worker
// @ID task-show
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=repositories.AgentModel}
// @Router /task/:id [get]
func (p *Task) Show(ctx *gin.Context) {
	userId := jwt.New().JwtUserId(ctx)
	taskIdUrl := ctx.Param("id")
	taskId, err := strconv.Atoi(taskIdUrl)
	if taskId == 0 || err != nil {
		response.BadRequest(ctx, "введите корректный task_id")
		return
	}

	task, err := repositories.TaskRepository().GetById(taskId)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}
	if task.UserID != userId {
		response.NotFound(ctx, "it's not your task")
		return
	}

	response.Data(ctx, task)
}
