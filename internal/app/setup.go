package app

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog/logger"
	"github.com/ut-proj/midiserver/pkg/types"
	"github.com/ut-proj/midiserver/pkg/version"
)

// Setup
func Setup(flags *types.Flags) {
	setupLogging(flags.LogLevel, flags.LogReportCaller)
	setupRandom()
	log.Info("Welcome to the Go midiserver!")
	log.Infof("running version: %s", version.VersionedBuildString())
	log.Infof("remote node: %s", flags.RemoteNode)
	log.Tracef("flags: %+v", flags)
}

// setupLogging ...
func setupLogging(logLevel string, reportCaller bool) {
	logger.SetupLogging(&logger.ZyLogOptions{
		Colored:      true,
		Level:        logLevel,
		Output:       "stderr",
		ReportCaller: reportCaller,
	})

}

// setupRandom ...
func setupRandom() {
	rand.Seed(time.Now().Unix())
}
