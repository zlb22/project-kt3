// +build appengine

package plog

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
