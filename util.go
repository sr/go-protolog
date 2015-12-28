package protolog

import (
	"fmt"
	"reflect"

	"go.pedge.io/protolog/pb"

	"github.com/golang/protobuf/proto"
)

// NOTE: the jsoonpb.Marshaler was EPICALLY SLOW in benchmarks
// When using the stdlib json.Marshal function instead for the text Marshaller,
// a speedup of 6X was observed!

func messageToEntryMessage(message proto.Message) (*protologpb.Entry_Message, error) {
	if message == nil {
		return nil, nil
	}
	value, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return &protologpb.Entry_Message{
		Name:  messageName(message),
		Value: value,
	}, nil
}

func entryMessageToMessage(entryMessage *protologpb.Entry_Message) (proto.Message, error) {
	if entryMessage == nil {
		return nil, nil
	}
	message, err := newMessage(entryMessage.Name)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(entryMessage.Value, message); err != nil {
		return nil, err
	}
	return message, nil
}

func messagesToEntryMessages(messages []proto.Message) ([]*protologpb.Entry_Message, error) {
	if messages == nil {
		return nil, nil
	}
	entryMessages := make([]*protologpb.Entry_Message, len(messages))
	for i, message := range messages {
		entryMessage, err := messageToEntryMessage(message)
		if err != nil {
			return nil, err
		}
		entryMessages[i] = entryMessage
	}
	return entryMessages, nil
}

func entryMessagesToMessages(entryMessages []*protologpb.Entry_Message) ([]proto.Message, error) {
	if entryMessages == nil {
		return nil, nil
	}
	messages := make([]proto.Message, len(entryMessages))
	for i, entryMessage := range entryMessages {
		message, err := entryMessageToMessage(entryMessage)
		if err != nil {
			return nil, err
		}
		messages[i] = message
	}
	return messages, nil
}

func newMessage(name string) (proto.Message, error) {
	reflectType := proto.MessageType(name)
	if reflectType == nil {
		return nil, fmt.Errorf("protolog: no Message registered for name: %s", name)
	}

	return reflect.New(reflectType.Elem()).Interface().(proto.Message), nil
}

func messageName(message proto.Message) string {
	if message == nil {
		return ""
	}
	return proto.MessageName(message)
}
