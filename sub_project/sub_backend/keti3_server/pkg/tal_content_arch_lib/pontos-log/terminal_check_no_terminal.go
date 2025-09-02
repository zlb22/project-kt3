// +build js nacl plan9

package plog

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
