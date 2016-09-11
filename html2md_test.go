package html2md

import (
	"testing"
	"bytes"
)

func TestBoldStrong(t *testing.T) {
	md, _ := HTML2MD(bytes.NewReader([]byte(`<b>a</b>`)))
	if *md != "**a**" && *md != "__a__" {
		t.Error("<b> should produce ** or __")
	}
	md, _ = HTML2MD(bytes.NewReader([]byte(`<strong>a</strong>`)))
	if *md != "**a**" && *md != "__a__" {
		t.Error("<strong> should produce ** or __")
	}
}

func TestIEm(t *testing.T) {
	md, _ := HTML2MD(bytes.NewReader([]byte(`<i>a</i>`)))
	if *md != "*a*" && *md != "_a_" {
		t.Error("<i> should produce * or _")
	}
	md, _ = HTML2MD(bytes.NewReader([]byte(`<em>a</em>`)))
	if *md != "*a*" && *md != "_a_" {
		t.Error("<i> should produce * or _")
	}
}

func TestConsecutiveBI(t *testing.T) {
	md, _ := HTML2MD(bytes.NewReader([]byte(`<b>a</b><i>b</i>`)))
	a, b := (*md)[0], (*md)[len(*md) -1]
	if (a == '*' && b == '_') || (a == '_' && b == '*') {
		t.Error("<b></b><i></i> should alternate formatting characters")
	}
}

func TestCode(t *testing.T) {
	md, _ := HTML2MD(bytes.NewReader([]byte(`<code>a</code>`)))
	if *md != "`a`" {
		t.Error("<code> should produce `")
	}
}
