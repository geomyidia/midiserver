package app

import (
	"math/rand"
	"time"

	"github.com/geomyidia/zylog/logger"
)

// SetupApp ...
func SetupLogging(logLevel string, reportCaller bool) {
	logger.SetupLogging(&logger.ZyLogOptions{
		Colored:      true,
		Level:        logLevel,
		Output:       "stderr",
		ReportCaller: reportCaller,
	})

}

func SetupRandom() {
	rand.Seed(time.Now().Unix())
}
