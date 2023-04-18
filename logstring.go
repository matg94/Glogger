package glogger

type LogString struct {
	str string
}

func Str(str string) *LogString {
	return &LogString{
		str: str,
	}
}

func (l *LogString) Error() string {
	return l.str
}
