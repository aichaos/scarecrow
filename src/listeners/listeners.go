/*
Package listeners provides implementations for various front-ends for the
Scarecrow chatbot.
*/
package listeners

import (
	"errors"
	"github.com/aichaos/scarecrow/src/types"
	"sort"
	"sync"
)

// Type Listener is an interface for front-end listeners for Scarecrow.
type Listener interface {
	New(types.ListenerConfig, chan types.CommunicationChannel, chan types.CommunicationChannel) Listener
	Start()
	InputChannel() chan types.CommunicationChannel
}

var (
	listenersMu sync.RWMutex
	listeners   = make(map[string]Listener)
)

// Register registers a listener handler.
func Register(name string, impl Listener) {
	listenersMu.Lock()
	defer listenersMu.Unlock()
	if impl == nil {
		panic("listeners: Registered listener is nil")
	}
	if _, dup := listeners[name]; dup {
		panic("listeners: Register called twice for listener " + name)
	}
	listeners[name] = impl
}

// Listeners returns a sorted list of the names of the registered listeners.
func Listeners() []string {
	listenersMu.RLock()
	defer listenersMu.RUnlock()
	var list []string
	for name := range listeners {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func Create(name string, cfg types.ListenerConfig, reply, answer chan types.CommunicationChannel) (Listener, error) {
	if _, ok := listeners[name]; !ok {
		return nil, errors.New("Unknown listener type.")
	}

	inst := listeners[name].New(cfg, reply, answer)
	return inst, nil
}
