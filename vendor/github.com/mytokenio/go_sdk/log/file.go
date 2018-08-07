package log

import (
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

type FileHandler struct {
	logger map[string]*logrus.Logger
}

func fileInit(conf *Config) *FileHandler {
	var (
		iw, ww, dw, ew *os.File
		err            error
		fh             = &FileHandler{
			logger: make(map[string]*logrus.Logger),
		}
	)
	ilg := logrus.New()
	iw, err = os.OpenFile(path.Join(conf.Dir, "info.log"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	ilg.Out = iw
	fh.logger["info"] = ilg

	wlg := logrus.New()
	ww, err = os.OpenFile(path.Join(conf.Dir, "warn.log"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	wlg.Out = ww
	fh.logger["warn"] = wlg

	dlg := logrus.New()
	dw, err = os.OpenFile(path.Join(conf.Dir, "debug.log"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	dlg.Out = dw
	fh.logger["debug"] = dlg

	elg := logrus.New()
	ew, err = os.OpenFile(path.Join(conf.Dir, "info.log"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	elg.Out = ew
	fh.logger["error"] = elg

	return fh
}

func (fh *FileHandler) info(format string, args ...interface{}) {
	fh.logger["info"].Infof(format, args...)
}

func (fh *FileHandler) warn(format string, args ...interface{}) {
	fh.logger["warn"].Warnf(format, args...)
}

func (fh *FileHandler) debug(format string, args ...interface{}) {
	fh.logger["debug"].Debugf(format, args...)
}

func (fh *FileHandler) error(format string, args ...interface{}) {
	fh.logger["error"].Errorf(format, args...)
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if s.IsDir() == false {
		panic("log path not dir")
	}
	return s.IsDir()
}

func isFile(path string) bool {
	return !isDir(path)
}
