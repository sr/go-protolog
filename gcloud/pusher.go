package gcloud

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"go.pedge.io/protolog"
	"google.golang.org/api/logging/v1beta3"
)

const customServiceName = "compute.googleapis.com"

var marshaler = &jsonpb.Marshaler{}

type pusher struct {
	service   *logging.ProjectsLogsEntriesService
	projectId string
	logName   string
}

func newPusher(client *http.Client, projectId string, logName string) *pusher {
	service, err := logging.New(client)
	if err != nil {
		panic(err)
	}
	return &pusher{service.Projects.Logs.Entries, projectId, logName}
}

func (p *pusher) Push(entry *protolog.Entry) error {
	logEntry, err := p.newLogEntry(entry)
	if err != nil {
		return err
	}
	request := p.service.Write(
		p.projectId,
		p.logName,
		&logging.WriteLogEntriesRequest{
			Entries: []*logging.LogEntry{logEntry},
		},
	)
	_, err = request.Do()
	return err
}

func (p *pusher) Flush() error {
	return nil
}

func (p *pusher) newLogEntry(entry *protolog.Entry) (*logging.LogEntry, error) {
	payload, err := p.marshalEntry(entry)
	if err != nil {
		return nil, err
	}

	return &logging.LogEntry{
		TextPayload: payload,
		Metadata: &logging.LogEntryMetadata{
			ServiceName: customServiceName,
		},
	}, nil
}

func (p *pusher) marshalEntry(entry *protolog.Entry) (string, error) {
	return marshaler.MarshalToString(entry)
}
