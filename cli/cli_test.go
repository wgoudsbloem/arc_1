package cli

import (
	"bytes"
	"os"
	"testing"
)

var (
	topic0 = "topic1"
	topic1 = "topic2"
	topic2 = "failtopic"
)

// use topic
func TestExecute(t *testing.T) {
	// select topic
	use := "use"
	cmdString := use + " " + topic0
	expMsg := topic0 + ">"
	cmd, msg := execute(cmdString)
	if msg != expMsg {
		t.Errorf("want: %v, got: %v", expMsg, msg)
	}
	// next test
	expMsg2 := `{"test":"value"}`
	cmdMsg2 := "put " + expMsg2
	msgPut := cmd(cmdMsg2)
	if msgPut[:5] != msgPutOk[:5] {
		t.Errorf("want: %v, got: %v", msgPutOk[:5], msgPut[:5])
	}
	// next test
	msgLast := cmd("last")
	if msgLast != expMsg2 {
		t.Errorf("want: %v, got: %v", expMsg2, msgLast)
	}
	// check for list command
	_, msgList := execute("list")
	expList := topic0
	if msgList != expList {
		t.Errorf("want:\n%v\ngot:\n%v", expList, msgList)
	}
}

func TestExecuteFail(t *testing.T) {
	// unkown command testing
	cmdNil := "dunno"
	if _, msgNil := execute(cmdNil); msgNil != msgCmdNotFound {
		t.Errorf("want: [%v], got: [%v]", msgCmdNotFound, msgNil)
	}
	// select topic
	use := "use"
	topic3 := ""
	cmdString := use + " " + topic3
	cmd, msg := execute(cmdString)
	if msg != msgNoTopicDef {
		t.Errorf("want: %v, got: %v", msgNoTopicDef, msg)
	}
	// setup valid topic for cmd testing"
	expMsg := topic2 + ">"
	cmdString = use + " " + topic2
	cmd, msg = execute(cmdString)
	if msg != expMsg {
		t.Errorf("want: %v, got: %v", expMsg, msg)
	}
	// ask for the last one should fail
	cmdMsg3 := "last"
	msgLast3 := cmd(cmdMsg3)
	if msgLast3 != msgTopicEmpty {
		t.Errorf("want: %v, got: %v", msgTopicEmpty, msgLast3)
	}
	// next test json is not valid
	expMsg2 := `{"test":"value`
	cmdMsg2 := "put " + expMsg2
	msgPut := cmd(cmdMsg2)
	if msgPut != msgPutFailInvalidJSON {
		t.Errorf("want: %v, got: %v", msgPutFailInvalidJSON, msgPut)
	}
	// next test should failt too, command is wrong
	msgLast := cmd("las")
	if msgLast != msgCmdNotFound {
		t.Errorf("want: %v, got: %v", msgCmdNotFound, msgLast)
	}
}

func TestUse(t *testing.T) {
	cmd, err := use(topic1)
	if err != nil {
		t.Error(err)
	}
	n := `put {"some":"json}`
	res := cmd(n)
	if res != msgPutFailInvalidJSON {
		t.Errorf("want: %v, got: %v", msgPutFailInvalidJSON, res)
	}
	//t.Log(res)
	x, err := use(topic1)
	if err != nil {
		t.Error(err)
	}
	exp0 := `{"some":"more"}`
	res2 := x("put " + exp0)
	if res2[:10] != msgPutOk[:10] {
		t.Errorf("want: %v, got: %v", msgPutOk, res2)
	}
	p := "last"
	res0 := cmd(p)
	if res0 != exp0 {
		t.Errorf("want: %v, got: %v", exp0, res0)
	}
}

func TestStart(t *testing.T) {
	cmd := "use"
	cmdString := cmd + " " + topic1
	expMsg := ">" + topic1 + ">\n"
	var r, w bytes.Buffer
	r.WriteString(cmdString)
	Start(&r, &w)
	msg, err := w.ReadString('\n')
	if err != nil {
		t.Error(err)
	}
	if msg != expMsg {
		t.Errorf("want: %v, got: %v", expMsg, msg)
	}
}

func TestCleanup(t *testing.T) {
	err := os.Remove(topic0 + ".topic")
	err = os.Remove(topic1 + ".topic")
	err = os.Remove(topic2 + ".topic")
	if err != nil {
		t.Error(err)
	}
}
