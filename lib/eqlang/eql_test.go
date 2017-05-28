package eqlang

import (
	"bytes"
	"testing"
)

// -----------------------------------------------------------------------------

const eqlTestCode = `<%
	//
	// eql example.eql -o example_bytes.go --imports=bytes --module=modbytes --Writer="*bytes.Buffer"
	// eql example.eql -o example_bufio.go --imports=bufio --module=modbufio --Writer="*bufio.Writer"
	//
%>
package eql_test

import (
	<%= eql.Imports() %>
	"encoding/binary"
)

// -----------------------------------------------------------------------------

type $module string

func (p $module) write(out $Writer, b []byte) {

	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

<% if Writer == "*bytes.Buffer" { %>
func (p $module) flush(out $Writer) {
}
<% } else { %>
func (p $module) flush(out $Writer) {

	err := out.Flush()
	if err != nil {
		panic(err)
	}
}
<% } %>

// -----------------------------------------------------------------------------
`

func TestEql(t *testing.T) {

	_, err := Parse(eqlTestCode)
	if err != nil {
		t.Fatal("Parse failed:", err)
	}
}

func TestSubst(t *testing.T) {

	out := Subst(`?$Writer!$`, map[string]interface{}{
		"Writer": "abc",
	})
	if out != "?abc!$" {
		t.Fatal("Subst failed:", out)
	}

	out = Subst(`$Writer!$$`, map[string]interface{}{
		"Writer": "abc",
	})
	if out != "abc!$" {
		t.Fatal("Subst failed:", out)
	}

	out = Subst(`$$$Writer!`, map[string]interface{}{
		"Writer": 123,
	})
	if out != "$123!" {
		t.Fatal("Subst failed:", out)
	}

	out = Subst(`$$$Writer`, map[string]interface{}{
		"Writer": 123,
	})
	if out != "$123" {
		t.Fatal("Subst failed:", out)
	}
}

func TestParseText(t *testing.T) {

	var b bytes.Buffer
	err := parseText(&b, "abc ``` def")
	if err != nil {
		t.Fatal("parseText failed:", err)
	}

	if b.String() != "print(eql.Subst(`abc ` + \"```\" + ` def`)); " {
		t.Fatal(b.String())
	}
}

// -----------------------------------------------------------------------------
