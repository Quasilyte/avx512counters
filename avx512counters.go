// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		{"validate command-line args", ctx.validateFlags},
		{"prepare work dir", ctx.prepareWorkDir},
		{"visit work dir", ctx.visitWorkDir},
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
	workDir       string
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
	flag.StringVar(&ctx.workDir, "workDir", "./avx512counters-workdir",
		`where to put results and the intermediate files`)
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

	absWorkDir, err := filepath.Abs(ctx.workDir)
	if err != nil {
		return fmt.Errorf("expand -workDir: %v", err)
	}
	ctx.workDir = absWorkDir
	return nil
}

func (ctx *context) validateFlags() error {
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
	if !fileExists(ctx.workDir) {
		if err := os.Mkdir(ctx.workDir, 0700); err != nil {
			return err
		}
	}

	// Always overwrite the main file, just in case.
	mainFile := filepath.Join(ctx.workDir, "main.go")
	mainFileContents := fmt.Sprintf(`
		// Code generated by avx512counters. DO NOT EDIT.
		package main
		func avx512routine(*[1024]byte)
		func main() {
			var memory [1024]byte
			for i := 0; i < %d; i++ {
				// Fill memory argument with some values.
				for i := range memory {
					memory[i] = byte(i)
				}
				avx512routine(&memory)
			}
		}`, ctx.loopCount)
	return ioutil.WriteFile(mainFile, []byte(mainFileContents), 0666)
}

func (ctx *context) visitWorkDir() error {
	return os.Chdir(ctx.workDir)
}

// fileExists reports whether file with given name exists.
func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
