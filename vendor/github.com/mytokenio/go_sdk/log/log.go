package log

import (
	"fmt"
	"time"
)

var (
	c  *Config
	tz *time.Location
	hs = make([]handle, 0)
)

type handle interface {
	info(format string, args ...interface{})
	warn(format string, args ...interface{})
	debug(format string, args ...interface{})
	error(format string, args ...interface{})
}

type xtime time.Time

// Config log config.
type Config struct {
	Stdout bool
	Dir    string
	Agent  *AgentConfig
}

// AgentConfig agent config.
type AgentConfig struct {
	TaskID  string
	Proto   string
	Addr    string
	Chan    int
	Timeout time.Duration
}

func (xt xtime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(xt).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// Init create logger with context.
func Init(conf *Config) {
	var (
		err error
	)
	tz, err = time.LoadLocation("Asia/shanghai")
	if err != nil {
		tz, _ = time.LoadLocation("PRC")
	}
	if conf.Dir != "" && isDir(conf.Dir) {
		hs = append(hs, fileInit(conf))
	}
	if conf.Agent != nil {
		hs = append(hs, agentInit(conf.Agent))
	}
	if conf.Stdout == true || len(hs) == 0 {
		hs = append(hs, stdInit())
	}
}

func Info(format string, args ...interface{}) {
	for _, h := range hs {
		h.info(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	for _, h := range hs {
		h.warn(format, args...)
	}
}
func Debug(format string, args ...interface{}) {
	for _, h := range hs {
		h.debug(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	for _, h := range hs {
		h.error(format, args...)
	}
}
