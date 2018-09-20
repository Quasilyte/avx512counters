// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"log"
	"regexp"
	"strings"
)

func main() {
	type step struct {
		name string
		fn   func() error
	}

	var ctx context

	steps := []step{
		{"init context", ctx.init},
		{"prepare work dir", ctx.prepareWorkDir},
	}

	for _, s := range steps {
		if err := s.fn(); err != nil {
			log.Fatalf("%s: %v", s.name, err)
		}
	}
}

type context struct {
	// memArgRE matches any kind of memory operand.
	// Displacement and indexing expressions are optional.
	memArgRE *regexp.Regexp

	// vmemArgRE is almost like memArgRE, but indexing expression is mandatory
	// and index register must be one of the X/Y/Z.
	vmemArgRE *regexp.Regexp

	// Fields below are initialized by command-line arguments (flags).

	extensions    []string
	perfTool      string
	iformSpanSize uint
	loopCount     uint
	perfRounds    uint
}

func (ctx *context) init() error {
	ctx.memArgRE = regexp.MustCompile(`(?:-?\d+)?\(\w+\)(?:\(\w+\*[1248]\))?`)
	ctx.vmemArgRE = regexp.MustCompile(`(?:-?\d+)?\(\w+\)\(([XYZ])\d+\*[1248]\)`)

	extensions := flag.String("extensions", "avx512f,avx512dq,avx512cd,avx512bw",
		`comma-separated list of extensions to be evaluated`)
	flag.StringVar(&ctx.perfTool, "perf", "perf",
		`perf tool binary name. ocperf and other drop-in replacements will do`)
	flag.UintVar(&ctx.iformSpanSize, "iformSpanSize", 100,
		`how many instruction lines form a single iform span. Higher values slow down the collection`)
	flag.UintVar(&ctx.loopCount, "loopCount", 1*1000*1000,
		`how many times to execute every iform span. Higher values slow down the collection`)
	flag.UintVar(&ctx.perfRounds, "perfRounds", 1,
		`how many times to re-validate perf results. Higher values slow down the collection`)

	flag.Parse()

	for _, ext := range strings.Split(*extensions, ",") {
		ext = strings.TrimSpace(ext)
		ctx.extensions = append(ctx.extensions, ext)
	}

	// Provide good error message to the user in case of issues.
	switch {
	case len(ctx.extensions) == 0:
		return errors.New("expected at least 1 extension name")
	case ctx.perfTool == "":
		return errors.New("argument -perf can't be empty")
	case ctx.iformSpanSize == 0:
		return errors.New("argument -iformSpanSize can't be 0")
	case ctx.loopCount == 0:
		return errors.New("argument -loopCount can't be 0")
	case ctx.perfRounds == 0:
		return errors.New("argument -perfRounds can't be 0")
	default:
		return nil
	}
}

func (ctx *context) prepareWorkDir() error {
	return nil
}
