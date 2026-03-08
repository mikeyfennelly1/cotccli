package libproducer

import (
	"context"
	"fmt"
	"time"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	log "github.com/sirupsen/logrus"
)

type ReaderDecorator struct {
	reader         Reader
	intervalSecs   time.Duration
	ctx            context.Context
	messageChannel chan client.Message
}

func (manager ReaderDecorator) GetName() string {
	return fmt.Sprintf("%s--%s", manager.reader.GetType(), manager.reader.GetName())
}

type Reader interface {
	GetType() string
	GetName() string
	GetValues() (map[string]float64, error)
	ToProducer() Producer
}

func (manager ReaderDecorator) StartScheduledProducer(collectorClient *client.CollectorClient, producer *client.CreatedProducer, streamName string) {
	ticker := time.NewTicker(manager.intervalSecs * time.Second)
	readerInstanceName := manager.reader.GetName()
	defer ticker.Stop()

	// goroutine for subscriber to this producer job
	go func() {
		for msg := range manager.messageChannel {
			log.Debugf("%s - sending message from producer to stream: %s", readerInstanceName, streamName)
			err := collectorClient.SendMessage(msg, streamName)
			if err != nil {
				log.Errorf("error writing to api: %v", err)
			}
		}
		log.Debugf("%s channel closed, exiting worker", manager.GetName())
	}()

	// start the producer
	for {
		select {
		case t := <-ticker.C:
			log.Infof("Scheduled reader job %s running at: %v", readerInstanceName, t)
			// get the values
			values, err := manager.reader.GetValues()
			if err != nil {
				log.Errorf("reader returned error when trying to get values: %v", err)
			}

			data := client.Message{
				ProducerId:   producer.UUID,
				ProducerName: readerInstanceName,
				ReadTime:     time.Now().UnixMilli(),
				Values:       values,
			}

			log.Tracef("writing message from reader %s to message channel", readerInstanceName)
			manager.messageChannel <- data

		case <-manager.ctx.Done():
			log.Infof("Scheduled reader task stopped for reader %s", readerInstanceName)
			return
		}
	}
}
