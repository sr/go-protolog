package protolog

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
	"unicode"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	jsonPBMarshaller = &jsonpb.Marshaler{}
)

type textMarshaller struct {
	options MarshallerOptions
}

func newTextMarshaller(options MarshallerOptions) *textMarshaller {
	return &textMarshaller{options}
}

func (t *textMarshaller) Marshal(goEntry *GoEntry) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if t.options.EnableID {
		_, _ = buffer.WriteString(goEntry.ID)
		_ = buffer.WriteByte(' ')
	}
	if !t.options.DisableTimestamp {
		_, _ = buffer.WriteString(goEntry.Time.Format(time.RFC3339))
		_ = buffer.WriteByte(' ')
	}
	if !t.options.DisableLevel {
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
			if err := t.marshalMessage(buffer, goEntry.Event); err != nil {
				return nil, err
			}
		}
	}
	if goEntry.Contexts != nil && len(goEntry.Contexts) > 0 && !t.options.DisableContexts {
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
				if err := t.marshalMessage(buffer, context); err != nil {
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

func (t *textMarshaller) marshalMessage(buffer *bytes.Buffer, message proto.Message) error {
	s, err := jsonPBMarshaller.MarshalToString(message)
	if err != nil {
		return err
	}
	_, _ = buffer.WriteString(messageName(message))
	_ = buffer.WriteByte(' ')
	_, _ = buffer.WriteString(s)
	return nil
}

func trimRightSpaceBytes(b []byte) []byte {
	return bytes.TrimRightFunc(b, unicode.IsSpace)
}
