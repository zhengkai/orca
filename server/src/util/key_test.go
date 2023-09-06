package util

import (
	"testing"
)

func TestKey(t *testing.T) {

	li := []string{
		`foo no-cache`,
		`no-cache foo`,
		`foo  no-cache`,
		`no-cache  foo`,
		` no-cache   foo`,
		` no-cache   foo `,
		`  foo    no-cache  `,
	}

	for _, s := range li {
		if !KeyNoCache(s) {
			t.Error(`KeyNoCache fail`, s)
		}
		if FormatKey(s) != `foo` {
			t.Error(`FormatKey fail`, s)
		}
	}

	li = []string{
		`foo`,
		`foono-cache`,
		`no-cachefoo`,
		` nocache foo`,
	}
	for _, s := range li {
		if KeyNoCache(s) {
			t.Error(`KeyNoCache fail`, s)
		}
	}
}
