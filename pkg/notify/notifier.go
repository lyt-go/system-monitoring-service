package notify

import "errors"

var ErrUnavailable = errors.New("告警通知器不可用")

type Notifier interface {
	Send(message string) error
}

type WebhookNotifier struct {
	Endpoint string
	SendFunc func(string) error
}

func (n *WebhookNotifier) Send(message string) error {
	if n.SendFunc != nil {
		return n.SendFunc(message)
	}
	if n.Endpoint == "" {
		return ErrUnavailable
	}
	return nil
}

func Available(notifier Notifier) bool {
	return notifier != nil
}
