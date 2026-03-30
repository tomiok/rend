package host

import (
	_ "embed"
	"encoding/json"
	"math"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

//go:embed host.html
var InfoHtml string

type Information struct {
	CPU        []Cpu
	CPUPercent float64
	Memory     Mem
	Disk       Disk
	General    General
	Net        Net
}

type Cpu struct {
	CPU        int32
	CPUPercent float64
	Model      string
	ModelName  string
	Vendor     string
}

type Mem struct {
	Active    uint64
	Available uint64
	Used      uint64
	Total     uint64
}

type Disk struct {
	Used    uint64
	Free    uint64
	Total   uint64
	Percent float64
}

type General struct {
	Hostname      string
	OS            string
	Platform      string
	Uptime        uint64
	Procs         uint32
	KernelVersion string
}

type Net struct {
	IOCounter string
}

func NewInformation() (Information, error) {
	memory, err := memoryInformation()
	if err != nil {
		return Information{}, err
	}

	diskInfo, err := diskInformation()
	if err != nil {
		return Information{}, err
	}

	cpuInfo, err := cpuInformation()
	if err != nil {
		return Information{}, err
	}

	general, err := generalInformation()
	if err != nil {
		return Information{}, err
	}

	return Information{
		Memory:     memory,
		CPUPercent: cpuInfo[0].CPUPercent,
		Disk:       diskInfo,
		CPU:        cpuInfo,
		General:    general,
	}, nil
}

func generalInformation() (General, error) {
	_info, err := host.Info()
	if err != nil {
		return General{}, err
	}

	s := _info.String()
	g := General{}
	if err = json.Unmarshal([]byte(s), &g); err != nil {
		return General{}, err
	}
	return g, nil
}

func diskInformation() (Disk, error) {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return Disk{}, err
	}

	free := diskInfo.Free
	used := diskInfo.Used
	total := diskInfo.Total

	return Disk{
		Used:    used,
		Free:    free,
		Total:   total,
		Percent: diskInfo.UsedPercent,
	}, nil
}

func memoryInformation() (Mem, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return Mem{}, err
	}

	return Mem{
		Active:    memStats.Active,
		Available: memStats.Available,
		Used:      memStats.Used,
		Total:     memStats.Total,
	}, nil
}

func cpuInformation() ([]Cpu, error) {
	infoStat, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	cpuInfo := make([]Cpu, len(infoStat))
	for j, i := range infoStat {
		_cpu := Cpu{
			CPU:        i.CPU,
			CPUPercent: percent[0],
			Model:      i.Model,
			ModelName:  i.ModelName,
			Vendor:     i.VendorID,
		}

		cpuInfo[j] = _cpu
	}

	return cpuInfo, nil
}

func toGB(bytes uint64) float64 {
	return math.Round(float64(bytes)/1024/1024/1024*10) / 10
}

func toMB(bytes uint64) float64 {
	return math.Round(float64(bytes)/1024/1024*10) / 10
}
