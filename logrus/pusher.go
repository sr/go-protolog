package protolog_logrus

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"unicode"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"go.pedge.io/protolog"
	"go.pedge.io/protolog/pb"
)

var (
	levelToLogrusLevel = map[protolog.Level]logrus.Level{
		protolog.LevelNone:  logrus.InfoLevel,
		protolog.LevelDebug: logrus.DebugLevel,
		protolog.LevelInfo:  logrus.InfoLevel,
		protolog.LevelWarn:  logrus.WarnLevel,
		protolog.LevelError: logrus.ErrorLevel,
		protolog.LevelFatal: logrus.FatalLevel,
		protolog.LevelPanic: logrus.PanicLevel,
	}
)

type pusher struct {
	logger  *logrus.Logger
	lock    *sync.Mutex
	options PusherOptions
}

func newPusher(options PusherOptions) *pusher {
	logger := logrus.New()
	if options.Out != nil {
		logger.Out = options.Out
	}
	if options.Hooks != nil && len(options.Hooks) > 0 {
		for _, hook := range options.Hooks {
			logger.Hooks.Add(hook)
		}
	}
	if options.Formatter != nil {
		logger.Formatter = options.Formatter
	}
	return &pusher{logger, &sync.Mutex{}, options}
}

func (p *pusher) Push(entry *protolog.Entry) error {
	logrusEntry, err := p.getLogrusEntry(entry)
	if err != nil {
		return err
	}
	return p.logLogrusEntry(logrusEntry)
}

type flusher interface {
	Flush() error
}

type syncer interface {
	Sync() error
}

func (p *pusher) Flush() error {
	if p.options.Out != nil {
		if syncer, ok := p.options.Out.(syncer); ok {
			return syncer.Sync()
		} else if flusher, ok := p.options.Out.(flusher); ok {
			return flusher.Flush()
		}
	}
	return nil
}

func (p *pusher) getLogrusEntry(entry *protolog.Entry) (*logrus.Entry, error) {
	logrusEntry := logrus.NewEntry(p.logger)
	logrusEntry.Time = entry.Time
	logrusEntry.Level = levelToLogrusLevel[entry.Level]

	if entry.ID != "" {
		logrusEntry.Data["_id"] = entry.ID
	}
	if !p.options.DisableContexts {
		for _, context := range entry.Contexts {
			if context == nil {
				continue
			}
			switch context.(type) {
			case *protologpb.Fields:
				for key, value := range context.(*protologpb.Fields).Value {
					if value != "" {
						logrusEntry.Data[key] = value
					}
				}
			default:
				if err := addProtoMessage(logrusEntry, context); err != nil {
					return nil, err
				}
			}
		}
	}
	if entry.Event != nil {
		switch entry.Event.(type) {
		case *protologpb.Event:
			logrusEntry.Message = trimRightSpace(entry.Event.(*protologpb.Event).Message)
		case *protologpb.WriterOutput:
			logrusEntry.Message = trimRightSpace(string(entry.Event.(*protologpb.WriterOutput).Value))
		default:
			logrusEntry.Data["_event"] = proto.MessageName(entry.Event)
			if err := addProtoMessage(logrusEntry, entry.Event); err != nil {
				return nil, err
			}
		}
	}
	return logrusEntry, nil
}

func (p *pusher) logLogrusEntry(entry *logrus.Entry) error {
	if err := entry.Logger.Hooks.Fire(entry.Level, entry); err != nil {
		return err
	}
	reader, err := entry.Reader()
	if err != nil {
		return err
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err = io.Copy(entry.Logger.Out, reader)
	return err
}

func addProtoMessage(logrusEntry *logrus.Entry, message proto.Message) error {
	m, err := getFieldsForProtoMessage(message)
	if err != nil {
		return err
	}
	for key, value := range m {
		logrusEntry.Data[key] = value
	}
	return nil
}

func getFieldsForProtoMessage(message proto.Message) (map[string]interface{}, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	if _, err := buffer.Write(data); err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, 0)
	if err := json.Unmarshal(buffer.Bytes(), &m); err != nil {
		return nil, err
	}
	n := make(map[string]interface{}, len(m))
	for key, value := range m {
		switch value.(type) {
		case map[string]interface{}:
			data, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			n[key] = string(data)
		default:
			n[key] = value
		}
	}
	return n, nil
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}
