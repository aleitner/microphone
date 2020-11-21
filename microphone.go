package microphone

import "github.com/faiface/beep"

func Init() error {
	return nil
}

func Terminate() error {
	return nil
}

func OpenDefaultStream(sampleRate beep.SampleRate, inputChannels int) (s *Streamer, format beep.Format, err error) {
	return nil, beep.Format{}, nil
}

type Streamer struct {
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
	return nil
}

func (s *Streamer) Start() error {
	return nil
}

func (s *Streamer) Stop() error {
	return nil
}