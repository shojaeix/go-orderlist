package olist

import "testing"

// the list.running must be true after run
func TestRun(t *testing.T) {
	list := OrderList{}
	list.run()
	if list.running != true {
		t.Error("list.running is false after run()")
	}
}
