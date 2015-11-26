/*
Package listeners provides implementations for various front-ends for the
Scarecrow chatbot.
*/
package listeners

import (
	"github.com/aichaos/scarecrow/src/types"
)

// TODO: the interface stuff isn't done properly, but until it is you can refer
// to the defined interface in this file to see what functions a listener should
// implement.

type Listener interface {
	Name() string
	New(types.ListenerConfig) *Listener
	Start()
}
