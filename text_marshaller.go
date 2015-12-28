package protolog

import (
	"bytes"
	"encoding/json"
	"time"
	"unicode"

	"go.pedge.io/protolog/pb"

	"github.com/fatih/color"
	"github.com/golang/protobuf/proto"
)

var (
	levelToColorString = map[Level]string{
		LevelNone:  color.BlueString(LevelNone.String()),
		LevelDebug: color.WhiteString(LevelDebug.String()),
		LevelInfo:  color.BlueString(LevelInfo.String()),
		LevelWarn:  color.YellowString(LevelWarn.String()),
		LevelError: color.RedString(LevelError.String()),
		LevelFatal: color.RedString(LevelFatal.String()),
		LevelPanic: color.RedString(LevelPanic.String()),
	}
)

type textMarshaller struct {
	disableTime     bool
	disableLevel    bool
	disableContexts bool
	colorize        bool
}

func newTextMarshaller(options ...TextMarshallerOption) *textMarshaller {
	textMarshaller := &textMarshaller{
		false,
		false,
		false,
		false,
	}
	for _, option := range options {
		option(textMarshaller)
	}
	return textMarshaller
}

func (t *textMarshaller) WithColors() TextMarshaller {
	return &textMarshaller{
		t.disableTime,
		t.disableLevel,
		t.disableContexts,
		true,
	}
}

func (t *textMarshaller) WithoutColors() TextMarshaller {
	return &textMarshaller{
		t.disableTime,
		t.disableLevel,
		t.disableContexts,
		false,
	}
}

func (t *textMarshaller) Marshal(entry *Entry) ([]byte, error) {
	return textMarshalEntry(entry, t.disableTime, t.disableLevel, t.disableContexts, t.colorize)
}

func textMarshalEntry(
	entry *Entry,
	disableTime bool,
	disableLevel bool,
	disableContexts bool,
	colorize bool,
) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if entry.ID != "" {
		_, _ = buffer.WriteString(entry.ID)
		_ = buffer.WriteByte(' ')
	}
	if !disableTime {
		_, _ = buffer.WriteString(entry.Time.Format(time.RFC3339))
		_ = buffer.WriteByte(' ')
	}
	if !disableLevel {
		var levelString string
		if colorize {
			levelString = levelToColorString[entry.Level]
		} else {
			levelString = entry.Level.String()
		}
		_, _ = buffer.WriteString(levelString)
		_ = buffer.WriteByte(' ')
		if len(levelString) == 4 {
			_ = buffer.WriteByte(' ')
		}
	}
	if entry.Event != nil {
		switch entry.Event.(type) {
		case *protologpb.Event:
			_, _ = buffer.WriteString(entry.Event.(*protologpb.Event).Message)
		case *protologpb.WriterOutput:
			_, _ = buffer.Write(trimRightSpaceBytes(entry.Event.(*protologpb.WriterOutput).Value))
		default:
			if err := textMarshalMessage(buffer, entry.Event); err != nil {
				return nil, err
			}
		}
	}
	if len(entry.Contexts) > 0 && !disableContexts {
		_ = buffer.WriteByte(' ')
		lenContexts := len(entry.Contexts)
		for i, context := range entry.Contexts {
			switch context.(type) {
			case *protologpb.Fields:
				data, err := json.Marshal(context.(*protologpb.Fields).Value)
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
				_ = buffer.WriteByte(' ')
			}
		}
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
