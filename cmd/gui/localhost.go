package gui

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// System Ststus
type LocalDuOShost struct {
	Cpu              []cpu.InfoStat                   `json:"cpu"`
	CpuPercent       []float64                        `json:"cpupercent"`
	Memory           mem.VirtualMemoryStat            `json:"mem"`
	Disk             disk.UsageStat                   `json:"disk"`
}

func (r *rcvar) GetlocalDuOShost() (l LocalDuOShost) {
	r.localhost = *new(LocalDuOShost)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	r.localhost.Cpu = sc
	r.localhost.CpuPercent = sp
	r.localhost.Memory = *sm
	r.localhost.Disk = *sd
	return
}
