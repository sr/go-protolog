package protolog

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
	"unicode"

	"github.com/golang/protobuf/proto"
)

var (
	defaultTextMarshallerOptions = textMarshallerOptions{}
)

type textMarshallerOptions struct {
	disableTime     bool
	disableLevel    bool
	disableContexts bool
}

type textMarshaller struct {
	options textMarshallerOptions
}

func newTextMarshaller(options ...TextMarshallerOption) *textMarshaller {
	textMarshallerOptions := textMarshallerOptions{
		false,
		false,
		false,
	}
	for _, option := range options {
		option(&textMarshallerOptions)
	}
	return &textMarshaller{textMarshallerOptions}
}

func (t *textMarshaller) Marshal(goEntry *GoEntry) ([]byte, error) {
	return textMarshalGoEntry(goEntry, t.options)
}

func textMarshalGoEntry(goEntry *GoEntry, options textMarshallerOptions) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if goEntry.ID != "" {
		_, _ = buffer.WriteString(goEntry.ID)
		_ = buffer.WriteByte(' ')
	}
	if !options.disableTime {
		_, _ = buffer.WriteString(goEntry.Time.Format(time.RFC3339))
		_ = buffer.WriteByte(' ')
	}
	if !options.disableLevel {
		levelString := strings.Replace(goEntry.Level.String(), "LEVEL_", "", -1)
		_, _ = buffer.WriteString(levelString)
		if len(levelString) == 4 {
			_, _ = buffer.WriteString("  ")
		} else {
			_ = buffer.WriteByte(' ')
		}
	}
	if goEntry.Event != nil {
		switch goEntry.Event.(type) {
		case *Event:
			_, _ = buffer.WriteString(goEntry.Event.(*Event).Message)
		case *WriterOutput:
			_, _ = buffer.Write(trimRightSpaceBytes(goEntry.Event.(*WriterOutput).Value))
		default:
			if err := textMarshalMessage(buffer, goEntry.Event); err != nil {
				return nil, err
			}
		}
	}
	if len(goEntry.Contexts) > 0 && !options.disableContexts {
		_, _ = buffer.WriteString(" contexts=[")
		lenContexts := len(goEntry.Contexts)
		for i, context := range goEntry.Contexts {
			switch context.(type) {
			case *Fields:
				data, err := json.Marshal(context.(*Fields).Value)
				if err != nil {
					return nil, err
				}
				_, _ = buffer.Write(data)
			default:
				if err := textMarshalMessage(buffer, context); err != nil {
					return nil, err
				}
			}
			if i != lenContexts-1 {
				_, _ = buffer.WriteString(", ")
			}
		}
		_ = buffer.WriteByte(']')
	}
	return trimRightSpaceBytes(buffer.Bytes()), nil
}

func textMarshalMessage(buffer *bytes.Buffer, message proto.Message) error {
	if message == nil {
		return nil
	}
	_, _ = buffer.WriteString(messageName(message))
	_ = buffer.WriteByte(' ')
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = buffer.Write(data)
	return err
}

func trimRightSpaceBytes(b []byte) []byte {
	return bytes.TrimRightFunc(b, unicode.IsSpace)
}
