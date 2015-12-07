package gcloud

import (
	"net/http"

	"go.pedge.io/protolog"
	"google.golang.org/api/logging/v1beta3"
)

const customServiceName = "compute.googleapis.com"

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
	request := p.service.Write(
		p.projectId,
		p.logName,
		&logging.WriteLogEntriesRequest{
			Entries: []*logging.LogEntry{
				&logging.LogEntry{
					TextPayload: "boomtown from protolog",
					Metadata: &logging.LogEntryMetadata{
						ServiceName: customServiceName,
					},
				},
			},
		},
	)
	_, err := request.Do()
	return err
}

func (p *pusher) Flush() error {
	return nil
}
