package libproducer

import (
	"context"
	"fmt"
	"net"

	"github.com/mikeyfennelly1/ise--y2--b3--project--desktop-sysinfo/client"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	log "github.com/sirupsen/logrus"
)

type sysinfoReader struct {
	id string
}

func (reader sysinfoReader) ToProducer() Producer {
	return ReaderDecorator{
		reader:         reader,
		intervalSecs:   1,
		ctx:            context.Background(),
		messageChannel: make(chan client.Message),
	}
}

func (reader sysinfoReader) GetType() string {
	return "sysinfo"
}

func (reader sysinfoReader) GetName() string {
	return fmt.Sprintf("sysinfo-reader--%s", reader.id)
}

func (reader sysinfoReader) GetValues() (map[string]float64, error) {
	log.Tracef("attempting to read virtual memory info")
	vmemStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	log.Tracef("attempting to read average load info")
	amount, err := load.Avg()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"vmem_available": float64(vmemStat.Available),
		"load_1":         amount.Load1,
		"load_5":         amount.Load5,
		"load_15":        amount.Load15,
	}, nil
}

func (reader sysinfoReader) getMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 ||
			len(iface.HardwareAddr) == 0 {
			continue
		}

		return iface.HardwareAddr.String(), nil
	}

	return "", fmt.Errorf("no MAC address found")
}
