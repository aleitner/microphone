package microphone

import (
	"github.com/faiface/beep"
	"github.com/gen2brain/malgo"
)

func OpenStream(ctx *malgo.AllocatedContext, deviceConfig malgo.DeviceConfig) (*Streamer, beep.Format, error) {
	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)
	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))

	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes
		newCapturedSampleCount := capturedSampleCount + sampleCount
		pCapturedSamples = append(pCapturedSamples, pSample...)
		capturedSampleCount = newCapturedSampleCount
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return nil, beep.Format{}, err
	}

	s := &Streamer{
		device: device,
	}

	format := beep.Format{
		SampleRate:  beep.SampleRate(device.SampleRate()),
		NumChannels: int(device.CaptureChannels()),
		Precision:   0,
	}

	return s, format, nil
}

type Streamer struct {
	device *malgo.Device
	buffer [][]float32
	err    error
}

func (s *Streamer) Stream(samples [][2]float64) (int, bool) {
	return 0, false
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