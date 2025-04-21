package orchestrator

import (
	"time"

	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/repositories"
	"github.com/grigoriy-st/YL-Golang/internal/orchestrator/services"
	"github.com/sirupsen/logrus"
)

func ResolveTasks() {
	repo := repositories.TaskRepository()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tasks, err := repo.GetTasksForResolve()
			if err != nil {
				logrus.Fatalf("Failed fetch tasks for resolve: %s", err.Error())
				break
			}

			for _, task := range tasks {
				if err := services.TaskService().ResolveTask(task); err != nil {
					logrus.Errorf("Failed send request to resolve task #%d: %s", task.TaskID, err.Error())
				}

				logrus.Infof("Resolve task #%d, expression \"%s\" has been added to queue", task.TaskID, task.Expression)
			}
		}
	}
}
