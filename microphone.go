package microphone

import (
	"fmt"
	"math"
	"sync"

	"github.com/faiface/beep"
	"github.com/gen2brain/malgo"
)

func OpenStream(ctx *malgo.AllocatedContext, deviceConfig malgo.DeviceConfig) (s *Streamer, format beep.Format, err error) {
	if deviceConfig.Capture.Channels > 2 || deviceConfig.Capture.Channels == 0 {
		return nil, beep.Format{}, fmt.Errorf("Invalid number of channels")
	}

	s = &Streamer{}

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))

	onRecvFrames := func(outputSample, inputSample []byte, framecount uint32) {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		samples := sampleBytesToFloats(inputSample, int(framecount), int(sizeInBytes), int(deviceConfig.Capture.Channels))
		s.buffer = append(s.buffer, samples...)
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return s, format, err
	}

	s.device = device

	format = beep.Format{
		SampleRate:  beep.SampleRate(device.SampleRate()),
		NumChannels: int(device.CaptureChannels()),
		Precision:   3,
	}

	return s, format, nil
}

type Streamer struct {
	mtx    sync.Mutex
	device *malgo.Device
	buffer [][2]float64
	err    error
	closed   bool
}

func (s *Streamer) Stream(samples [][2]float64) (int, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.closed {
		return 0, false
	}

	// Stream is already empty
	if !s.device.IsStarted() || len(s.buffer) == 0 {
		return 0, true
	}

	numSamples := len(samples)
	if len(s.buffer) < numSamples {
		numSamples = len(s.buffer)
	}

	numSamplesStreamed := 0
	for i := 0; i < numSamples; i++ {
		if len(s.buffer) == 0 {
			break
		}

		samples[i] = s.buffer[i]
		numSamplesStreamed++
	}

	s.buffer = s.buffer[numSamplesStreamed:]

	return numSamplesStreamed, true
}

func (s *Streamer) Err() error {
	return s.err
}

func (s *Streamer) Close() error {
	s.Stop()

	if !s.closed {
		s.device.Uninit()
		s.closed = true
	}
	return nil
}

func (s *Streamer) Start() {
	if !s.device.IsStarted() {
		s.device.Start()
	}
}

func (s *Streamer) Stop() {
	if s.device.IsStarted() {
		s.device.Stop()
	}
}

func sampleBytesToFloats(input []byte, sampleCount, sampleSizeInBytes, numChannels int) [][2]float64 {
	samples := make([][2]float64, sampleCount)

	if numChannels == 0 || numChannels > 2 {
		return samples
	}

	for i := range samples {
		for channel := 0; channel < numChannels; channel++ {
			bytes := input[:sampleSizeInBytes]
			samples[i][channel] = decodeFloat(bytes)
			input = input[sampleSizeInBytes:]
		}
	}

	return samples
}

func decodeFloat(p []byte) (x float64) {
	precision := 3

	var xUint64 uint64
	for i := precision - 1; i >= 0; i-- {
		xUint64 <<= 8
		xUint64 += uint64(p[i])
	}

	if xUint64 >= 1<<uint(precision*8-1) {
		compl := 1<<uint(precision*8) - xUint64
		return -float64(int64(compl)) / (math.Exp2(float64(precision)*8-1) - 1)
	}
	return float64(int64(xUint64)) / (math.Exp2(float64(precision)*8-1) - 1)
}
