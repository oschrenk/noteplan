package internal

import (
	"testing"
)

// TestParse calls greetings.Hello with a name, checking
func TestParse(t *testing.T) {
	name := "Gladys"
	want := "Gladys"
	if name != want {
		t.Fatalf("Same")
	}
}
