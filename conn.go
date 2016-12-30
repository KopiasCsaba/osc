package osc

import (
	"context"
	"errors"
	"net"
	"strings"
)

const (
	bufSize = 4096 // Size of read and write buffers.
)

// Common errors.
var (
	ErrNilDispatcher  = errors.New("nil dispatcher")
	ErrPrematureClose = errors.New("server cannot be closed before calling Listen")
)

// Conn defines the methods
type Conn interface {
	net.Conn

	Context() context.Context
	Serve(Dispatcher) error
	Send(Packet) error
	SendTo(net.Addr, Packet) error
}

var invalidAddressRunes = []rune{'*', '?', ',', '[', ']', '{', '}', '#', ' '}

// ValidateAddress returns an error if addr contains
// characters that are disallowed by the OSC spec.
func ValidateAddress(addr string) error {
	for _, chr := range invalidAddressRunes {
		if strings.ContainsRune(addr, chr) {
			return ErrInvalidAddress
		}
	}
	return nil
}
