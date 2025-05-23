package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/grigoriy-st/YL-Golang/pkg/rabbitmq"
	"time"
)

func Ping(queueOrchestrator string, agentId string, pingTime int) {
	ticker := time.NewTicker(time.Duration(pingTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			answer := amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(agentId),
				Type:        "ping",
				Timestamp:   time.Now(),
			}

			if err := rabbitmq.Get().SendToQueue(queueOrchestrator, answer); err != nil {
				logrus.Fatalf("Failed sent ping: %s", err.Error())
				break
			}

			logrus.Debugf("Ping was successful sent")
		}
	}
}
