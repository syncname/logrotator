package logrotator

type StdWriter struct {
	rotator *LogRotator
}

func (cw *StdWriter) Write(data []byte) (int, error) {
	return cw.rotator.Write(data)
}

func (cw *StdWriter) SetRotator(rotator *LogRotator) {
	cw.rotator = rotator
}
