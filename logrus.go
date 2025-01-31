package logrotator

type LogrusWriter struct {
	rotator *LogRotator
}

func (cw *LogrusWriter) Write(data []byte) (int, error) {
	return cw.rotator.Write(data)
}

func (cw *LogrusWriter) SetRotator(rotator *LogRotator) {
	cw.rotator = rotator
}
