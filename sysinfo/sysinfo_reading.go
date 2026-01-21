package sysinfo

import (
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetReading() (map[string]float64, error) {
	vmemStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
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
