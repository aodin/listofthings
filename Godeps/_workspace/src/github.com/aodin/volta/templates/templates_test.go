package templates

import (
	"bytes"
	"net/http"
	"testing"
)

// testWriter implements the net/http ResponseWriter interface
// Write is implemented by the embedded bytes.Buffer
type testWriter struct {
	*bytes.Buffer
}

// Header returns a new header everytime
func (tw *testWriter) Header() http.Header {
	return http.Header{}
}

// WriteHeader does nothing with the code
func (tw *testWriter) WriteHeader(code int) {}

func newTestWriter() *testWriter {
	var b bytes.Buffer
	return &testWriter{&b}
}

// Check that testWriter implements the http.ResponseWriter
var _ http.ResponseWriter = newTestWriter()

func TestEmpty(t *testing.T) {
	empty := Empty()
	if err := empty.Add("{{fsfd"); err == nil {
		t.Errorf("templates: Add failed to error when given a bad template")
	}
	if err := empty.Add(`{{ define "whatever" }}{{ .Whatever }}{{ end }}`); err != nil {
		t.Errorf("templates: Add returned an error when given a proper template")
	}

	// The whatever template should be added to the templates
	var b []byte
	buffer := bytes.NewBuffer(b)
	empty.Execute(buffer, "whatever")
}

func TestTemplates(t *testing.T) {
	// Parse the test_fixtures directory with a local variable
	parsed := New("./test_fixtures/pass", Attrs{"Greeting": "Yo"})

	// Overwrite a local variable
	if err := parsed.SetAttr("Greeting", "Hello"); err == nil {
		t.Errorf("failed to warn of overwritten local attributes")
	}

	// Create a testWriter for testing output
	w := newTestWriter()

	// Render the template with an additional Name attribute
	parsed.Execute(w, "parent", Attrs{"Name": "Lebowski"})

	// Test the output versus expected. If you've modified the test_fixtures
	// templates watch out for newlines - they matter!
	expected := `Parent is named Lebowski
Child says Hello`
	if w.String() != expected {
		t.Errorf("unexpected template execution output: %s", w.String())
	}

	// Test alternate delims
	parsed = NewWithDelims(
		"./test_fixtures/delims",
		`<%`,
		`%>`,
		Attrs{"Greeting": "Yo"},
	)
	w = newTestWriter()
	parsed.Execute(w, "delims")
	expected = `The template says Yo`
	if w.String() != expected {
		t.Errorf("unexpected delim template execution output: %s", w.String())
	}

	// Call a template that doesn't exist
	// Execute the test in an anon block (in case we want more defers)
	func() {
		var panicked interface{}
		defer func() {
			panicked = recover()
		}()
		w = newTestWriter()
		parsed.Execute(w, "dne")
		if panicked == nil {
			t.Fatalf("failed to panic when given a name that does not exist")
		}
	}()

	// Cause an error in filepath.Walk by providing a dir that does not exist
	func() {
		var panicked interface{}
		defer func() {
			panicked = recover()
		}()
		New("./test_fixtures/dne")
		if panicked == nil {
			t.Fatalf("failed to panic when given a dir that does not exist")
		}
	}()

	// Parse a bad template
	func() {
		var panicked interface{}
		defer func() {
			panicked = recover()
		}()
		New("./test_fixtures/fail")
		if panicked == nil {
			t.Fatalf("failed to panic when a bad template is given")
		}
	}()
}
