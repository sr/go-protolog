package protolog_syslog

import (
	"log/syslog"

	"go.pedge.io/protolog"
)

var (
	levelToLogFunc = map[protolog.Level]func(*syslog.Writer, string) error{
		protolog.Level_LEVEL_NONE:  (*syslog.Writer).Info,
		protolog.Level_LEVEL_DEBUG: (*syslog.Writer).Debug,
		protolog.Level_LEVEL_INFO:  (*syslog.Writer).Info,
		protolog.Level_LEVEL_WARN:  (*syslog.Writer).Warning,
		protolog.Level_LEVEL_ERROR: (*syslog.Writer).Err,
		protolog.Level_LEVEL_FATAL: (*syslog.Writer).Crit,
		protolog.Level_LEVEL_PANIC: (*syslog.Writer).Alert,
	}
)

type pusher struct {
	writer     *syslog.Writer
	marshaller protolog.Marshaller
}

func newPusher(writer *syslog.Writer, options ...PusherOption) *pusher {
	pusher := &pusher{writer, DefaultTextMarshaller}
	for _, option := range options {
		option(pusher)
	}
	return pusher
}

func (p *pusher) Flush() error {
	return nil
}

func (p *pusher) Push(goEntry *protolog.GoEntry) error {
	data, err := p.marshaller.Marshal(goEntry)
	if err != nil {
		return err
	}
	return levelToLogFunc[goEntry.Level](p.writer, string(data))
}
