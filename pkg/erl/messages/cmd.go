package messages

import (
	"github.com/ut-proj/midiserver/pkg/erl/datatypes"
)

const (
	ArgsKey    string = "args"
	CommandKey string = "command"
)

type Cmd struct {
	commandType *datatypes.Atom
	args        *datatypes.List
}

func NewCmd(tuples ...datatypes.Tuple) *Cmd {
	return &Cmd{}
}

func (c *Cmd) Type() string {
	return c.commandType.Value()
}

func (c *Cmd) Args() *datatypes.List {
	return c.args
}

func (c *Cmd) ToTerm() {

}
