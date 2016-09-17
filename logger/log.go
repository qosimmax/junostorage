package logger

import (
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
)

var (
	log  *logrus.Logger
	once sync.Once
)

func GetLogger() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()
		//Output to stderr
		log.Out = os.Stderr
	})

	return log
}
