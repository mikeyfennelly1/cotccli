package libproducer

import (
	"fmt"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	log "github.com/sirupsen/logrus"
)

type Producer interface {
	StartScheduledProducer(collectorClient *client.CollectorClient, producer *client.CreatedProducer, streamName string)
	GetName() string
}

func ReaderFactory(sourceType string, sourceId string) (Reader, error) {
	switch sourceType {
	case "sysinfo":
		return sysinfoReader{
			id: sourceId,
		}, nil
	default:
		log.Fatalf("reader type %s does not exist", sourceType)
		return nil, fmt.Errorf("unknown reader type - provided: '%s'", sourceType)
	}
}
