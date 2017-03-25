package pubsub

import (
	"fmt"
	"testing"
)

type TestMsg struct {
	val string
}

var ps PubSuber

func init() {
	ps = PubSuber{}
}

func TestSubscribe(t *testing.T) {
	testMsg := TestMsg{"test message!"}
	fn := func(msg interface{}) error {
		_msg := msg.(TestMsg)
		if _msg.val != testMsg.val {
			return fmt.Errorf("expected '%s', but got '%s'", _msg.val, testMsg.val)
		}
		return nil
	}
	ps.Subscribe(fn)
}

func TestNotify(t *testing.T) {
	testMsg := TestMsg{"test message!"}
	ps.Notify(testMsg)
}
