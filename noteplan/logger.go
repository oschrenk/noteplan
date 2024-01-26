package noteplan

import (
	"log"
	"os"
)

type Log struct {
	Enabled bool
}

var mylog = log.New(os.Stderr, "", log.LstdFlags)

var Logger = Log{
	Enabled: false,
}

func (logger *Log) Log(s string) {
	if logger.Enabled {
		mylog.Println(s)
	}
}
