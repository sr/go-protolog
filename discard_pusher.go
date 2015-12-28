package protolog

var (
	discardPusherInstance = &discardPusher{}
)

type discardPusher struct{}

func (d *discardPusher) Flush() error {
	return nil
}

func (d *discardPusher) Push(_ *GoEntry) error {
	return nil
}
