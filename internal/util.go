package kacker

import (
	"os/exec"
)

func hasCommand(cmd string) bool {
	ret := exec.Command(cmd)
	if err := ret.Run(); err != nil {
		return false
	}
	return true
}

// HasPacker checks whether packer is installed on the system.
func HasPacker() bool {
	return hasCommand("packer")
}

// HasKsvalidator checks whether ksvalidator, as part of pykickstart, is
// installed on the system.
func HasKsvalidator() bool {
	return hasCommand("ksvalidator")
}
