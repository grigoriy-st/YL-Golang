package agents

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/repositories"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/services"
	"strconv"
)

func HandleAnswer(message amqp.Delivery) {
	// handle answer from agent
	taskId, err := strconv.Atoi(message.CorrelationId)
	if err != nil {
		logrus.Errorf("Get wrong task id for answer: %s", message.CorrelationId)
		return
	}

	// set answer into database
	if err = services.TaskService().SetAnswer(taskId, string(message.Body), repositories.STATUS_COMPLETED); err != nil {
		logrus.Errorf("Failed update a row with task %d: %s", taskId, err.Error())
		return
	}

	logrus.Infof("Get answer for %s: %s", message.CorrelationId, message.Body)
}
