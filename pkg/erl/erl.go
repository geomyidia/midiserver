package erl

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	NodeDelimit   = "@"
	ShortNodename = "midiserver"
)

var (
	LongNodename = fmt.Sprintf("%s%slocalhost", ShortNodename, NodeDelimit)
)

func ReadCookie() (string, error) {
	usr, _ := user.Current()
	cookie, err := os.ReadFile(filepath.Join(usr.HomeDir, ".erlang.cookie"))
	if err != nil {
		return "", err
	}
	return string(cookie), nil
}
