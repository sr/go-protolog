package protolog

func init() {
	Register("protolog.Fields", func() Message { return &Fields{} })
	Register("protolog.Event", func() Message { return &Event{} })
	Register("protolog.WriterOutput", func() Message { return &WriterOutput{} })
}

func (m *Fields) ProtologName() string {
	return "protolog.Fields"
}
func (m *Event) ProtologName() string {
	return "protolog.Event"
}
func (m *WriterOutput) ProtologName() string {
	return "protolog.WriterOutput"
}
