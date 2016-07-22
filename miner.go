package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/decred/gominer/cl"
)

func getCLPlatforms() ([]cl.CL_platform_id, error) {
	var numPlatforms cl.CL_uint
	status := cl.CLGetPlatformIDs(0, nil, &numPlatforms)
	if status != cl.CL_SUCCESS {
		return nil, clError(status, "CLGetPlatformIDs")
	}
	platforms := make([]cl.CL_platform_id, numPlatforms)
	status = cl.CLGetPlatformIDs(numPlatforms, platforms, nil)
	if status != cl.CL_SUCCESS {
		return nil, clError(status, "CLGetPlatformIDs")
	}
	return platforms, nil
}

// getCLDevices returns the list of devices for the given platform.
func getCLDevices(platform cl.CL_platform_id) ([]cl.CL_device_id, error) {
	var numDevices cl.CL_uint
	status := cl.CLGetDeviceIDs(platform, cl.CL_DEVICE_TYPE_GPU, 0, nil,
		&numDevices)
	if status != cl.CL_SUCCESS {
		return nil, clError(status, "CLGetDeviceIDs")
	}
	devices := make([]cl.CL_device_id, numDevices)
	status = cl.CLGetDeviceIDs(platform, cl.CL_DEVICE_TYPE_ALL, numDevices,
		devices, nil)
	if status != cl.CL_SUCCESS {
		return nil, clError(status, "CLGetDeviceIDs")
	}
	return devices, nil
}

type Miner struct {
	devices          []*Device
	workDone         chan []byte
	quit             chan struct{}
	needsWorkRefresh chan struct{}
	wg               sync.WaitGroup
	pool             *Stratum

	validShares   uint64
	invalidShares uint64
}

func NewMiner() (*Miner, error) {
	m := &Miner{
		workDone:         make(chan []byte, 10),
		quit:             make(chan struct{}),
		needsWorkRefresh: make(chan struct{}),
	}

	// If needed, start pool code.
	if cfg.Pool != "" && !cfg.Benchmark {
		s, err := StratumConn(cfg.Pool, cfg.PoolUser, cfg.PoolPassword)
		if err != nil {
			return nil, err
		}
		m.pool = s
	}

	platformIDs, err := getCLPlatforms()
	if err != nil {
		return nil, fmt.Errorf("Could not get CL platforms: %v", err)
	}
	platformID := platformIDs[0]
	deviceIDs, err := getCLDevices(platformID)
	if err != nil {
		return nil, fmt.Errorf("Could not get CL devices for platform: %v", err)
	}

	m.devices = make([]*Device, len(deviceIDs))
	for i, deviceID := range deviceIDs {
		var err error
		m.devices[i], err = NewDevice(i, platformID, deviceID, m.workDone)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Miner) workSubmitThread() {
	defer m.wg.Done()

	for {
		select {
		case <-m.quit:
			return
		case data := <-m.workDone:
			// Only use that is we are not using a pool.
			if m.pool == nil {
				accepted, err := GetWorkSubmit(data)
				if err != nil {
					m.invalidShares++
					minrLog.Errorf("Error submitting work: %v", err)
				} else {
					if accepted {
						m.validShares++
					} else {
						m.invalidShares++
					}
					minrLog.Errorf("Submitted work successfully: %v", accepted)
					m.needsWorkRefresh <- struct{}{}
				}
			} else {
				accepted, err := GetPoolWorkSubmit(data, m.pool)
				if err != nil {
					m.invalidShares++
					minrLog.Errorf("Error submitting work to pool: %v", err)
				} else {
					if accepted {
						m.validShares++
					} else {
						m.invalidShares++
					}
					minrLog.Errorf("Submitted work to pool successfully: %v",
						accepted)
					m.needsWorkRefresh <- struct{}{}
				}
			}
		}
	}
}

func (m *Miner) workRefreshThread() {
	defer m.wg.Done()

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		// Only use that is we are not using a pool.
		if m.pool == nil {
			work, err := GetWork()
			if err != nil {
				minrLog.Errorf("Error in getwork: %v", err)
			} else {
				for _, d := range m.devices {
					d.SetWork(work)
				}
			}
		} else {
			work, err := GetPoolWork(m.pool)
			if err != nil {
				minrLog.Errorf("Error in getpoolwork: %v", err)
			} else {
				for _, d := range m.devices {
					d.SetWork(work)
				}
			}
		}
		select {
		case <-m.quit:
			return
		case <-t.C:
		case <-m.needsWorkRefresh:
		}
	}
}

func (m *Miner) printStatsThread() {
	defer m.wg.Done()

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()

	for {
		minrLog.Infof("Total accepted shares: %v, total rejected shares: %v",
			m.validShares, m.invalidShares)
		for _, d := range m.devices {
			d.PrintStats()
		}

		select {
		case <-m.quit:
			return
		case <-t.C:
		case <-m.needsWorkRefresh:
		}
	}
}

func (m *Miner) Run() {
	m.wg.Add(len(m.devices))

	for _, d := range m.devices {
		device := d
		go func() {
			device.Run()
			device.Release()
			m.wg.Done()
		}()
	}

	m.wg.Add(1)
	go m.workSubmitThread()

	if cfg.Benchmark {
		minrLog.Warn("Running in BENCHMARK mode! No real mining taking place!")
		work := &Work{}
		for _, d := range m.devices {
			d.SetWork(work)
		}
	} else {
		m.wg.Add(1)
		go m.workRefreshThread()
	}

	m.wg.Add(1)
	go m.printStatsThread()

	m.wg.Wait()
}

func (m *Miner) Stop() {
	close(m.quit)
	for _, d := range m.devices {
		d.Stop()
		m.wg.Done()
	}
}
