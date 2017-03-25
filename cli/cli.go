package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"arcessio/stream"
)

const (
	cmdUse  = "use"
	cmdList = "list"
	cmdPut  = "put"
	cmdLast = "last"
	cmdQuit = "end"
)

const (
	msgPutOk              = "success offset: %d"
	msgPutFailInvalidJSON = "JSON is not valid"
	msgCmdNotFound        = "command not found"
	msgNoTopicDef         = "no topic definded (usage: use <<topic>>)"
	msgTopicEmpty         = "topic has no values"
	msgPrmpt              = "%s>"
)

type command func(n string) string

var workDir = "./"

// Start is the cli entry point
func Start(r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)
	var cmd command
	var msg string
	w.Write([]byte(">"))
	for s.Scan() {
		if cmd != nil {
			msg = cmd(s.Text())
		} else {
			cmd, msg = execute(s.Text())
		}
		w.Write([]byte(msg + "\n>"))
	}
}

func execute(cmdStr string) (cmd command, msg string) {
	input := strings.Split(cmdStr, " ")
	var err error
	switch input[0] {
	default:
		msg = msgCmdNotFound
	case cmdQuit:
		os.Exit(0)
	case cmdUse:
		if len(input) < 2 || input[1] == "" {
			msg = msgNoTopicDef
			return
		}
		cmd, err = use(input[1])
		if err != nil {
			msg = err.Error()
		}
		msg = fmt.Sprintf(msgPrmpt, input[1])
	case cmdList:
		fs, err := ioutil.ReadDir(workDir)
		if err != nil {
			msg = err.Error()
		}
		for i, f := range fs {
			if strings.HasSuffix(f.Name(), ".topic") {
				msg += strings.TrimSuffix(f.Name(), ".topic")
				if i < len(fs)-1 {
					msg += "\n"
				}
			}
		}
	}
	return
}

func use(topic string) (cmd command, err error) {
	//needs check
	fs, err := stream.NewFileStreamReaderWriter(topic)
	if err != nil {
		return
	}
	cmd = func(cmdStr string) (s string) {
		n := strings.Split(cmdStr, " ")
		switch n[0] {
		case cmdPut:
			if !isJSON([]byte(n[1])) {
				s = msgPutFailInvalidJSON
				return
			}
			offset, err := fs.WriteByte([]byte(n[1]))
			if err != nil {
				return err.Error()
			}
			return fmt.Sprintf(msgPutOk, offset)
		case cmdLast:
			jsn, _, err := fs.LastJSON()
			if err != nil {
				switch t := err.(type) {
				case (*json.SyntaxError):
					if t.Error() == "unexpected end of JSON input" {
						return msgTopicEmpty
					}
				default:
					return err.Error()
				}
			}
			return string(jsn)
		case cmdQuit:
			os.Exit(0)
		default:
			return msgCmdNotFound
		}
		return
	}
	return
}

func cmdNotFound(s []string) command {
	return func(n string) string {
		return msgCmdNotFound
	}
}

func isJSON(jsn []byte) bool {
	var m map[string]interface{}
	return json.Unmarshal(jsn, &m) == nil
}
