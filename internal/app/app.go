package app

import (
	"math/rand"
	"time"

	"github.com/geomyidia/zylog/logger"
)

// SetupApp ...
func SetupLogging() {
	logger.SetupLogging(&logger.ZyLogOptions{
		Colored:      true,
		Level:        "debug",
		Output:       "stderr",
		ReportCaller: true,
	})

}

func SetupRandom() {
	rand.Seed(time.Now().Unix())
}
