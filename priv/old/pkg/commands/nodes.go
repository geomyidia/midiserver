package commands

import (
	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/epmd"
	"github.com/ut-proj/midiserver/pkg/types"
)

func ListNodes(flags *types.Flags) {
	nodes, err := epmd.ListNodes(flags.EPMDHost, flags.EPMDPort)
	if err != nil {
		log.Error(err)
	}
	for _, node := range nodes {
		println(node)
	}
}
