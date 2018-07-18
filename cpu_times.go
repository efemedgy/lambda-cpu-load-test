package cpuload

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
)

var pid string

func init() {
	pid = strconv.Itoa(os.Getpid())
}

type CpuTimesStatComparer struct {
	Before *CpuTimesStat
	After  *CpuTimesStat
}

type CpuTimesStat struct {
	*ProcPidTimesStat
	*ProcTimesStat
}

type ProcPidTimesStat struct {
	Utime uint64 `json:"utime"`
	Stime uint64 `json:"stime"`
}

type ProcTimesStat struct {
	User       uint64 `json:"user"`
	System     uint64 `json:"system"`
	Idle       uint64 `json:"idle"`
	Nice       uint64 `json:"nice"`
	Iowait     uint64 `json:"iowait"`
	Irq        uint64 `json:"irq"`
	SoftIrq    uint64 `json:"softirq"`
	Steal      uint64 `json:"steal"`
	Guest      uint64 `json:"guest"`
	Guest_Nice uint64 `json:"guest_nice"`
}

func (t *CpuTimesStat) total() uint64 {
	return t.User + t.Nice + t.System + t.Idle + t.Iowait
}

func (t *CpuTimesStat) sys_used() uint64 {
	return t.User + t.Nice + t.System
}

func (t *CpuTimesStat) proc_used() uint64 {
	return t.Utime + t.Stime
}

func Sample() *CpuTimesStat {
	pps := getProcPidStat(pid)
	ps := getProcStat()
	return &CpuTimesStat{
		ProcPidTimesStat: pps,
		ProcTimesStat:    ps,
	}
}

// Reads utime and stime from /proc/[pid]/stat file
func getProcPidStat(pid string) *ProcPidTimesStat {
	contents, err := ioutil.ReadFile("/proc/" + pid + "/stat")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fields := strings.Fields(string(contents))
	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		fmt.Println("procStat[13]: ", err.Error())
	}
	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		fmt.Println("procStat[13]: ", err.Error())
	}
	return &ProcPidTimesStat{
		Utime: utime,
		Stime: stime,
	}
}

// Reads stats from /proc/stat file
func getProcStat() *ProcTimesStat {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fields := strings.Fields(string(contents))
	user, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		fmt.Println("procStat[1] ", err.Error())
	}
	nice, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		fmt.Println("procStat[2] ", err.Error())
	}
	system, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		fmt.Println("procStat[3] ", err.Error())
	}
	idle, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		fmt.Println("procStat[4] ", err.Error())
	}
	iowait, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		fmt.Println("procStat[5] ", err.Error())
	}
	irq, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		fmt.Println("procStat[6] ", err.Error())
	}
	softirq, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		fmt.Println("procStat[7] ", err.Error())
	}
	steal, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		fmt.Println("procStat[8] ", err.Error())
	}
	guest, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		fmt.Println("procStat[9] ", err.Error())
	}
	guest_nice, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		fmt.Println("procStat[10] ", err.Error())
	}
	return &ProcTimesStat{
		User:       user,
		System:     system,
		Idle:       idle,
		Nice:       nice,
		Iowait:     iowait,
		Irq:        irq,
		SoftIrq:    softirq,
		Steal:      steal,
		Guest:      guest,
		Guest_Nice: guest_nice,
	}
}

func CalculateSystemUsagePercent(comparer *CpuTimesStatComparer) float64 {
	dSysUsed := comparer.After.sys_used() - comparer.Before.sys_used()
	dTotal := comparer.After.total() - comparer.Before.total()
	s := float64(dSysUsed) / float64(dTotal)
	if s <= 0 {
		s = 0
	} else if s >= 1 {
		s = 1
	}
	return s
}

func CalculateProcessUsagePercent(comparer *CpuTimesStatComparer) float64 {
	dProcUsed := comparer.After.proc_used() - comparer.Before.proc_used()
	dTotal := comparer.After.total() - comparer.Before.total()
	p := float64(dProcUsed) / float64(dTotal)
	if p <= 0 {
		p = 0
	} else if p >= 1 {
		p = 1
	}
	return p
}
