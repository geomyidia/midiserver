package app

import (
	"math/rand"
	"time"

	"github.com/geomyidia/zylog/logger"
)

// SetupApp ...
func SetupLogging(logLevel string) {
	logger.SetupLogging(&logger.ZyLogOptions{
		Colored:      true,
		Level:        logLevel,
		Output:       "stderr",
		ReportCaller: true,
	})

}

func SetupRandom() {
	rand.Seed(time.Now().Unix())
}
