package broker

type Publisher interface {
	Publish(event string, payLoad string)
}
