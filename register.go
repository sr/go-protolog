package protolog

import "fmt"

var (
	nameToMessageConstructor = make(map[string]func() Message)
	registeredMessages       = make(map[registeredMessage]bool)
)

type registeredMessage struct {
	name        string
	messageType MessageType
}

// Register registers a Message constructor funcation to a message name.
// This should only be called by generated code.
func Register(name string, messageType MessageType, messageConstructor func() Message) {
	registeredMessage := registeredMessage{name: name, messageType: messageType}
	if _, ok := registeredMessages[registeredMessage]; ok {
		panic(fmt.Sprintf("protolog: duplicate Message registered: %s", name))
	}
	registeredMessages[registeredMessage] = true
	nameToMessageConstructor[name] = messageConstructor
}
