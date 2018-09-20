// Copyright 2018 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestInstructionForm(t *testing.T) {
	tests := []struct {
		op    string
		args  []string
		iform string
	}{
		{
			"KANDW",
			[]string{"K4", "K4", "K6"},
			"KANDW K, K, K",
		},
		{
			"KMOVW",
			[]string{"K5", "-17(BP)(SI*4)"},
			"KMOVW K, mem",
		},
		{
			"VCMPPD",
			[]string{"$65", "X9", "X7", "K4", "K4"},
			"VCMPPD imm, X, X, K, K",
		},
		{
			"VCMPPD",
			[]string{"$0", "-17(BP)(SI*2)", "Z0", "K5", "K6"},
			"VCMPPD imm, mem, Z, K, K",
		},
		{
			"VCMPPS",
			[]string{"$81", "99(R15)(R15*2)", "Y16", "K4", "K1"},
			"VCMPPS imm, mem, Y, K, K",
		},
		{
			"VCVTSD2USIQ",
			[]string{"(CX)", "R13"},
			"VCVTSD2USIQ mem, reg",
		},
	}

	for i, test := range tests {
		line := testLine{op: test.op, args: test.args}
		have := instructionForm(line)
		want := test.iform
		if have != want {
			t.Errorf("[%d]: iforms mismatch:\nhave: %q\nwant: %q", i, have, want)
		}
	}
}

func TestScanner(t *testing.T) {
	filename := filepath.Join("testdata", "asmfile.s")
	scanner := testFileScanner{}
	if err := scanner.init(filename); err != nil {
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
