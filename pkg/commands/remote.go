package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ut-proj/midiserver/pkg/erl/epmd"
	"github.com/ut-proj/midiserver/pkg/erl/rpc"
	"github.com/ut-proj/midiserver/pkg/types"
)

func ShowRemotePort(flags *types.Flags) {
	port, err := epmd.NodePort(flags.EPMDHost, flags.EPMDPort, flags.RemoteNode)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("Remote node %s is running on port %d\n", flags.RemoteNode, port)
}

func PingRemoteModule(flags *types.Flags) error {
	rpcClient, err := rpc.New(flags)
	if err != nil {
		return err
	}
	result, err := rpcClient.Ping(flags.RemoteModule)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
