// Copyright (c) 2016 The Decred developers.

package main

import (
	"math"
	"time"
	"unsafe"

	"github.com/decred/gominer/cl"
	"github.com/decred/gominer/util"
	"github.com/decred/gominer/work"
)

// getKernelExecutionTime returns the kernel execution time for a device.
func (d *Device) getKernelExecutionTime(globalWorksize uint32) (time.Duration,
	error) {
	d.work = work.Work{}

	minrLog.Tracef("Started DEV #%d: %s for kernel execution time fetch",
		d.index, d.deviceName)
	outputData := make([]uint32, outputBufferSize)

	var status cl.CL_int

	// Calculate the precalculation for the first round optimization.
	var workArray [180]byte
	copy(workArray[:], d.work.Data[0:180])
	work32 := util.ConvertByteSliceHeaderToUint32Slice(workArray)
	h, v, xorLUT := precalculateStatesAndLUT(d.midstate, work32)
	hV, vV := *h, *v

	// arg 0: pointer to the buffer
	obuf := d.outputBuffer
	argument := 0
	status = cl.CLSetKernelArg(d.kernel, cl.CL_uint(argument),
		cl.CL_size_t(unsafe.Sizeof(obuf)),
		unsafe.Pointer(&obuf))
	if status != cl.CL_SUCCESS {
		return 0, clError(status, "CLSetKernelArg (output buffer)")
	}
	argument++

	// args 1..16: precomputed v
	for i := 0; i < 16; i++ {
		vi := vV[i]
		status = cl.CLSetKernelArg(d.kernel, cl.CL_uint(argument),
			uint32Size, unsafe.Pointer(&vi))
		if status != cl.CL_SUCCESS {
			return 0, clError(status, "CLSetKernelArg (v)")
		}
		argument++
	}

	// arg 17: last uint32 of midstate
	h1 := hV[1]
	status = cl.CLSetKernelArg(d.kernel, cl.CL_uint(argument),
		uint32Size, unsafe.Pointer(&h1))
	if status != cl.CL_SUCCESS {
		return 0, clError(status, "CLSetKernelArg (midstate)")
	}
	argument++

	// arg 18: the XOR precomputation LUT
	lutSize := uint32Size * 215
	cl_xorLUT := cl.CLCreateBuffer(d.context, cl.CL_MEM_READ_ONLY|
		cl.CL_MEM_COPY_HOST_PTR, lutSize, unsafe.Pointer(xorLUT), &status)
	if status != cl.CL_SUCCESS {
		return 0, clError(status, "CLCreateBuffer (xorLUT)")
	}
	status = cl.CLSetKernelArg(d.kernel, cl.CL_uint(argument),
		cl.CL_size_t(unsafe.Sizeof(cl_xorLUT)),
		unsafe.Pointer(&cl_xorLUT))
	if status != cl.CL_SUCCESS {
		return 0, clError(status, "CLSetKernelArg (xorLUT)")
	}
	argument++

	// Clear the found count from the buffer
	status = cl.CLEnqueueWriteBuffer(d.queue, d.outputBuffer,
		cl.CL_FALSE, 0, uint32Size, unsafe.Pointer(&zeroSlice[0]),
		0, nil, nil)
	if status != cl.CL_SUCCESS {
		return time.Duration(0), clError(status, "CLEnqueueWriteBuffer")
	}

	// Execute the kernel and follow its execution time.
	currentTime := time.Now()
	var globalWorkSize [1]cl.CL_size_t
	globalWorkSize[0] = cl.CL_size_t(globalWorksize)
	var localWorkSize [1]cl.CL_size_t
	localWorkSize[0] = localWorksize
	status = cl.CLEnqueueNDRangeKernel(d.queue, d.kernel, 1, nil,
		globalWorkSize[:], localWorkSize[:], 0, nil, nil)
	if status != cl.CL_SUCCESS {
		return time.Duration(0), clError(status, "CLEnqueueNDRangeKernel")
	}

	// Read the output buffer.
	cl.CLEnqueueReadBuffer(d.queue, d.outputBuffer, cl.CL_TRUE, 0,
		uint32Size*outputBufferSize, unsafe.Pointer(&outputData[0]), 0,
		nil, nil)
	if status != cl.CL_SUCCESS {
		return time.Duration(0), clError(status, "CLEnqueueReadBuffer")
	}

	// Release the local buffer for the LUT.
	cl.CLReleaseMemObject(cl_xorLUT)

	elapsedTime := time.Since(currentTime)
	minrLog.Tracef("DEV #%d: Kernel execution to read time for work "+
		"size calibration: %v", d.index, elapsedTime)

	return elapsedTime, nil
}

// calcWorkSizeForMilliseconds calculates the correct worksize to achieve
// a device execution cycle of the passed duration in milliseconds.
func (d *Device) calcWorkSizeForMilliseconds(ms int) (uint32, error) {
	workSize := uint32(1 << 10)
	timeToAchieve := time.Duration(ms) * time.Millisecond
	for {
		execTime, err := d.getKernelExecutionTime(workSize)
		if err != nil {
			return 0, err
		}

		// If we fail to go above the desired execution time, double
		// the work size and try again.
		if execTime < timeToAchieve {
			workSize <<= 1
			continue
		}

		// We're passed the desired execution time, so now calculate
		// what the ideal work size should be.
		adj := float64(workSize) * (float64(timeToAchieve) / float64(execTime))
		adj /= 256.0
		adjMultiple256 := uint32(math.Ceil(adj))
		workSize = adjMultiple256 * 256

		break
	}

	return workSize, nil
}
