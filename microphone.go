package microphone

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