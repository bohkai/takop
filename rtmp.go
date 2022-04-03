package main

import(
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

type Rtmp struct {
	client *rtmp.ClientConn
	stream *rtmp.Stream
}

func NewRtmp(url string) (*Rtmp, error) {
	client, err := rtmp.Dial("rtmp", url, &rtmp.ConnConfig{})
	if err != nil {
		return nil, err
	}
	defer client.Close()

	stream, err := client.CreateStream(nil)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	return &Rtmp{
		client: client,
		stream: stream,
	}, nil
}

func (r *Rtmp)Publish(name string) error {
	if err := r.stream.Publish(&rtmpmsg.NetStreamPublish{
		PublishingName: name,
		PublishingType: "live",
	}); err != nil {
		return err
	}
	return nil
}

func (r *Rtmp)Connect() error {
	if err := r.client.Connect(nil); err != nil {
		return err
	}

	return nil
}

func (r *Rtmp)Close() error {
	if err := r.client.Close(); err != nil {
		return err
	}

	return nil
}