package natsutils

import (
	"errors"
	"reflect"

	"github.com/nats-io/nats.go"
)

func Subscribe(nc *nats.Conn, subject string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nc.Subscribe(subject, func(msg *nats.Msg) {
		go cb(msg)
	})
}

func QueueSubscribe(nc *nats.Conn, subject, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		go cb(msg)
	})
}

func SubscribeEc(ec *nats.EncodedConn, subject string, cb nats.Handler) (*nats.Subscription, error) {
	wrappedCB, err := makeConcurrentHandler(cb)
	if err != nil {
		return nil, err
	}
	return ec.Subscribe(subject, wrappedCB)
}

func QueueSubscribeEc(ec *nats.EncodedConn, subject, queue string, cb nats.Handler) (*nats.Subscription, error) {
	wrappedCB, err := makeConcurrentHandler(cb)
	if err != nil {
		return nil, err
	}
	return ec.QueueSubscribe(subject, queue, wrappedCB)
}

func makeConcurrentHandler(cb nats.Handler) (nats.Handler, error) {
	rf := reflect.TypeOf(cb)
	if rf.Kind() != reflect.Func {
		return nil, errors.New("nats: Handler should be a function")
	}

	vf := reflect.ValueOf(cb)

	return reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		go vf.Call(in)
		return nil
	}).Interface(), nil
}

func NewEncodedJsonConn(url string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return ec, nil
}

func TopicParse(serviceName string, topics []string) string {
	subBase := serviceName
	for _, topic := range topics {
		subBase = subBase + "." + topic
	}

	return subBase
}
