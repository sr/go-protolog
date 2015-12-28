/*
Package protolog_syslog defines functionality for integration with syslog.
*/
package protolog_syslog // import "go.pedge.io/protolog/syslog"

import (
	"log/syslog"

	"go.pedge.io/protolog"
)

var (
	// DefaultTextMarshaller is the default text Marshaller for syslog.
	DefaultTextMarshaller = protolog.NewTextMarshaller(
		protolog.TextMarshallerDisableTime(),
		protolog.TextMarshallerDisableLevel(),
	)
)

// PusherOptions defines options for constructing a new syslog protolog.Pusher.
type PusherOptions struct {
	// By default, DefaultTextMarshaller is used.
	Marshaller protolog.Marshaller
}

// NewPusher creates a new protolog.Pusher that logs using syslog.
func NewPusher(writer *syslog.Writer, options PusherOptions) protolog.Pusher {
	return newPusher(writer, options)
}
