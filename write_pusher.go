package protolog

import (
	"bytes"
	"io"
	"sync"
)

var (
	newlineBytes = []byte{'\n'}
)

type writePusher struct {
	writer     io.Writer
	marshaller Marshaller
	newline    bool
	lock       *sync.Mutex
}

func newWritePusher(writer io.Writer, options WritePusherOptions) *writePusher {
	writePusher := &writePusher{
		writer,
		options.Marshaller,
		options.Newline,
		&sync.Mutex{},
	}
	if writePusher.marshaller == nil {
		writePusher.marshaller = DelimitedMarshaller
	}
	return writePusher
}

type flusher interface {
	Flush() error
}

type syncer interface {
	Sync() error
}

func (w *writePusher) Flush() error {
	if syncer, ok := w.writer.(syncer); ok {
		return syncer.Sync()
	} else if flusher, ok := w.writer.(flusher); ok {
		return flusher.Flush()
	}
	return nil
}

func (w *writePusher) Push(goEntry *GoEntry) error {
	data, err := w.marshaller.Marshal(goEntry)
	if err != nil {
		return err
	}
	if w.newline {
		buffer := bytes.NewBuffer(data)
		_, _ = buffer.Write(newlineBytes)
		data = buffer.Bytes()
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	if _, err := w.writer.Write(data); err != nil {
		return err
	}
	return nil
}
