package microphone

import (
	"sync"

	"github.com/faiface/beep"
	"github.com/gen2brain/malgo"
)

func OpenStream(ctx *malgo.AllocatedContext, deviceConfig malgo.DeviceConfig) (*Streamer, beep.Format, error) {
	s := &Streamer{
		minBufferSize: 4096,
	}

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))

	onRecvFrames := func(outputSample, inputSample []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes

		samples := sampleBytesToFloats(inputSample, int(sampleCount))
		s.buffer = append(s.buffer, samples...)
		if len(s.buffer) > s.minBufferSize {
			s.stop()
		}
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return nil, beep.Format{}, err
	}

	s.device = device

	format := beep.Format{
		SampleRate:  beep.SampleRate(device.SampleRate()),
		NumChannels: int(device.CaptureChannels()),
		Precision:   3,
	}

	s.start()

	return s, format, nil
}

type Streamer struct {
	mtx    sync.Mutex
	device *malgo.Device
	buffer [][2]float64
	minBufferSize int
	err    error
}

func (s *Streamer) Stream(samples [][2]float64) (int, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// Stream is already empty
	if len(s.buffer) == 0 {
		return 0, false
	}

	numSamplesStreamed := 0
	for i := range samples {
		if len(s.buffer) == 0 {
			break
		}

		samples[i] = s.buffer[i]
		numSamplesStreamed++
	}

	s.buffer = s.buffer[numSamplesStreamed:]

	if len(s.buffer) < s.minBufferSize {
		s.start()
	}

	return numSamplesStreamed, true
}

func (s *Streamer) Err() error {
	return s.err
}

func (s *Streamer) Close() error {
	s.device.Uninit()

	return nil
}

func (s *Streamer) start() {
	s.device.Start()
}

func (s *Streamer) stop() {
	s.device.Stop()
}

func sampleBytesToFloats(input []byte, sampleCount int) [][2]float64 {
	samples := make([][2]float64, sampleCount)

	return samples
}