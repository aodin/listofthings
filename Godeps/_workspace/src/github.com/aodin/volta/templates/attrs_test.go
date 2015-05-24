package templates

import (
	"html/template"
	"testing"
)

type testUser struct {
	Name string `json:"name"`
}

func TestAttrs(t *testing.T) {
	// Create a new attrs map
	a := Attrs{"Greeting": "Hello"}

	// And another
	b := Attrs{"Greeting": "Yo", "ID": 1}

	// Calling the method merge on a should overwrite any duplicate values
	// on a with those in b
	a.Merge(b)

	// The greeting should be overwritten - remember it is an interface{}!
	greetingAttr := a["Greeting"]
	greeting, ok := greetingAttr.(string)
	if !ok {
		t.Fatalf("The Greeting attr could not be cast to a string")
	}
	if greeting != "Yo" {
		t.Fatalf("The Greeting attr was not Yo")
	}

	// And the ID should be added
	IDAttr := a["ID"]
	ID, ok := IDAttr.(int)
	if !ok {
		t.Fatalf("The ID attr could not be cast to an int")
	}
	if ID != 1 {
		t.Fatalf("The ID attr was not 1")
	}

	attrs := AsJSON("user", testUser{Name: "admin"})
	j, ok := attrs["user"]
	if !ok {
		t.Fatalf("The json attr has no user key")
	}
	if string(j.(template.JS)) != `{"name":"admin"}` {
		t.Fatalf("incorrect json output: %s", j)
	}

	// The following type cannot be marshaled
	invalid := AsJSON("user", map[int64]bool{1: true})
	if _, ok := invalid["user"]; ok {
		t.Fatalf("an attr was created from invalid JSON")
	}
}
