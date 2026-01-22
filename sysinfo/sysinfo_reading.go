package sysinfo

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	log "github.com/sirupsen/logrus"
)

func GetReading() (map[string]float64, error) {
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

	sinceEpochMilli := time.Now().UnixMilli()

	return map[string]float64{
		"read_time":      float64(sinceEpochMilli),
		"vmem_available": float64(vmemStat.Available),
		"load_1":         amount.Load1,
		"load_5":         amount.Load5,
		"load_15":        amount.Load15,
	}, nil
}

func ScheduledProducer(ctx context.Context, messageChannel chan map[string]float64) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			log.Infof("Scheduled sysinfo job running at: %v", t)
			data, err := GetReading()
			if err != nil {
				log.Errorf("error reading sysinfo: %v", err)
			}

			log.Tracef("writing sysinfo to sysinfo channel")
			messageChannel <- data

		case <-ctx.Done():
			log.Infof("Scheduled sysinfo task stopped")
			return
		}
	}
}
