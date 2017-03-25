package stream

import (
	"testing"
)

func TestInternalLastEntry(t *testing.T) {
	expVal1 := `{"test1":"val1"}`
	expVal2 := `{"test2":"val2"}`
	testVal1 := expVal1 + "\n"
	testVal2 := expVal2 + "\n"
	testVal3 := testVal1 + testVal2
	res1 := lastEntry([]byte(testVal1))
	if string(res1) != expVal1 {
		t.Errorf("want: '%v' got: '%v'", expVal1, string(res1))
	}
	res2 := lastEntry([]byte(testVal3))
	if string(res2) != expVal2 {
		t.Errorf("want: '%v' got: '%v'", expVal2, string(res2))
	}
}
