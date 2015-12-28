package protolog

import (
	"bytes"
	"io"

	"go.pedge.io/protolog/pb"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
)

type delimitedMarshaller struct{}

func (m *delimitedMarshaller) Marshal(goEntry *GoEntry) ([]byte, error) {
	entry, err := goEntry.ToEntry()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	if _, err := pbutil.WriteDelimited(buffer, entry); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

type delimitedUnmarshaller struct{}

func (u *delimitedUnmarshaller) Unmarshal(reader io.Reader, goEntry *GoEntry) error {
	entry := &protologpb.Entry{}
	if _, err := pbutil.ReadDelimited(reader, entry); err != nil {
		return err
	}
	iGoEntry, err := EntryToGoEntry(entry)
	if err != nil {
		return err
	}
	*goEntry = *iGoEntry
	return nil
}
