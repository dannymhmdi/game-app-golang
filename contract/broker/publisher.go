package broker

type Publisher interface {
	Publish(event string, payLoad string) error
}

type PublisherFunc func(event string, payLoad string) error
