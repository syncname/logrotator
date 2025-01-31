package logrotator

type SlogWriter struct {
	rotator *LogRotator
}

func (cw *SlogWriter) Write(data []byte) (int, error) {
	return cw.rotator.Write(data)
}

func (cw *SlogWriter) SetRotator(rotator *LogRotator) {
	cw.rotator = rotator
}
