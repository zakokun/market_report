package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

type stdHandler struct {
	logger *logrus.Logger
}

func stdInit() *stdHandler {
	sh := new(stdHandler)
	lg := logrus.New()
	lg.Out = os.Stdout
	sh.logger = lg
	return sh
}

func (fh *stdHandler) info(format string, args ...interface{}) {
	fh.logger.Infof(format, args...)
}

func (fh *stdHandler) warn(format string, args ...interface{}) {
	fh.logger.Warnf(format, args...)
}

func (fh *stdHandler) debug(format string, args ...interface{}) {
	fh.logger.Debugf(format, args...)
}

func (fh *stdHandler) error(format string, args ...interface{}) {
	fh.logger.Errorf(format, args...)
}
