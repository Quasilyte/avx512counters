package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	scanner := testFileScanner{filename: filepath.Join("testdata", "asmfile.s")}
	if err := scanner.init(); err != nil {
		t.Fatalf("scanner init error: %v", err)
	}
	expected := []testLine{
		{
			"VAESDEC",
			[]string{"X24", "X7", "X11"},
			"\tVAESDEC X24, X7, X11                               // 62124508ded8",
		},
		{
			"VAESENCLAST",
			[]string{"7(SI)(DI*1)", "Z6", "Z11"},
			"\tVAESENCLAST 7(SI)(DI*1), Z6, Z11                   // 62724d48dd9c3e07000000",
		},
	}
	for i := 0; scanner.scan(); i++ {
		have := scanner.line
		want := expected[i]
		if have.op != want.op {
			t.Errorf("[%d]: op mismatch:\nhave: %q\nwant: %q",
				i, have.op, want.op)
		}
		if !reflect.DeepEqual(have.args, want.args) {
			t.Errorf("[%d]: args mismatch:\nhave: %v\nwant: %v",
				i, have.args, want.args)
		}
		if have.text != want.text {
			t.Errorf("[%d]: text mismatch:\nhave: %q\nwant: %q",
				i, have.text, want.text)
		}
	}
}
