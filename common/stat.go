package common

import (
	"runtime"

	"sam-api/models"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func NewStat() *models.Stat {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s := &models.Stat{
		m.Alloc,
		m.TotalAlloc,
		m.Sys,
		m.NumGC,
	}
	return s
}
