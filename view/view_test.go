package view3

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"arcessio/pubsub"
	"arcessio/stream"
)

var (
	expectedID    = "id123"
	expectedTopic = "customer"
	expectedTotal float64
)

var customers = []string{
	`{"name":"john", "acct":100}`,
	`{"name":"john", "acct":-50}`,
	`{"name":"john", "acct":25}`}

func TestCreateCustomerTopic(t *testing.T) {
	sw, err := stream.NewFileStreamWriter("customer")
	if err != nil {
		t.Fatal(err)
	}
	//	for i := 0; i < len(customers); i++ {
	//		sw.WriteJson([]byte(customers[i]))
	//	}
	for i := 0; i < 100; i++ {
		f := math.Ceil(rand.Float64() * 100)
		if i%2 == 0 {
			f = f * -1
		}
		expectedTotal = expectedTotal + f
		sw.WriteByteArray([]byte(fmt.Sprintf(`{"name":"john", "second":"mnbvcxzlkjhgfdsapoiuytrewq", "acct":%v}`, f)))
	}
}

//Need to be able to subscribe for notification when view has been build
//passed json schema needs to be unmarshalled to Schema struct
//

func TestNewViewer(t *testing.T) {
	//pass a schema and a subsriber to start the Viewer build process
	schema := []byte(fmt.Sprintf(`{"id":"%v","topic":"%s","sum":"acct"}`, expectedID, expectedTopic))
	var subscriber pubsub.Subscriber
	subscriber = func(view interface{}) error {
		v, ok := view.([]byte)
		if !ok {
			t.Errorf("expected view to be View")
		}
		testView(t, v)
		return nil
	}
	viewer, err := NewViewer(schema, subscriber)
	if err != nil {
		t.Fatal(err)
	}
	if viewer == nil {
		t.Fatal("expect viewer not to be nil")
	}
	if viewer.schema.ID != expectedID {
		t.Errorf("expected id to be %v, but got %v", expectedID, viewer.schema.ID)
	}
	if viewer.schema.Topic != expectedTopic {
		t.Errorf("expected topic to be %v, but got %v", expectedTopic, viewer.schema.Topic)
	}
}

func testView(t *testing.T, v []byte) {
	t.Log("subtest: testView")
	var result struct {
		Name   string  `json:"name"`
		Acct   float64 `json:"acct"`
		Offset int64
	}
	if err := json.Unmarshal(v, &result); err != nil {
		t.Error(err)
	}
	if result.Name != "john" {
		t.Errorf("expected name to be %v, but got %v", "john", result.Name)
	}
	if result.Acct != expectedTotal {
		t.Errorf("expected acct to be %v, but got %v", expectedTotal, result.Acct)
	}
}

func TestCleanup(t *testing.T) {
	os.Remove(expectedID + ".schema")
	os.Remove(expectedTopic + ".topic")
}
