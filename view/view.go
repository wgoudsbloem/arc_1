package view3

import (
	"encoding/json"
	"io"
	"os"

	"arcessio/pubsub"
	"arcessio/stream"
)

// Schema is structure
type Schema struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
	Sum   string `json:"sum"`
	Avg   string `json:"avg"`
}

// Viewer contains the properties that make up the view
type Viewer struct {
	pubsub.PubSub
	schema *Schema
	view   []byte
}

// NewViewer is a factory to get Viewer
func NewViewer(schema []byte, n pubsub.Subscriber) (viewer *Viewer, err error) {
	// marshal byte array into Schema struct
	var s Schema
	err = json.Unmarshal(schema, &s)
	if err != nil {
		return
	}
	// create the file and store the schema
	f, err := os.OpenFile(s.ID+".schema", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	_, err = f.Write(schema)
	// create the pubsub package and subscribe the passed Subscriber
	ps := &pubsub.PubSuber{}
	ps.Subscribe(n)
	// finally create the viewer and set the members
	// view is nil as it needs to be build
	viewer = &Viewer{ps, &s, nil}
	viewer.buildView()
	return
}

func (v *Viewer) buildView() (err error) {
	var viewMap map[string]interface{}
	//get the stream reader from the topic
	sr, err := stream.NewFileStreamReader(v.schema.Topic)
	if err != nil {
		return
	}
	var sum float64
	var jsn []byte
	var offset, offset0 int64
	for {
		offset0, err = sr.Read(&jsn)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		offset = offset0
		//log.Info(string(jsn))
		err = json.Unmarshal(jsn, &viewMap)
		if err != nil {
			return
		}
		if v.schema.Sum != "" {
			s, _ := viewMap[v.schema.Sum].(float64)
			sum = sum + s
		}
	}
	viewMap[v.schema.Sum] = sum
	viewMap["offset"] = offset
	v.view, err = json.Marshal(viewMap)
	if err != nil {
		return
	}
	v.Notify(v.view)
	return
}
