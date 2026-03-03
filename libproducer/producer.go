package libproducer

import log "github.com/sirupsen/logrus"

type Producer interface {
	StartScheduledProducer()
}

func ReaderFactory(sourceType string, sourceId string) Reader {
	switch sourceType {
	case "sysinfo":
		return sysinfoReader{
			id: sourceId,
		}
	default:
		log.Fatalf("reader type %s does not exist", sourceType)
		return nil
	}
}
