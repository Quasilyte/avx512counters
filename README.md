# avx512counters

AVX-512 hardware counters collector written in Go, based on Go toolchain.

## Overview

This program utilized Go 1.11 assembler AVX-512 support, extensive end2end test suite
and Linux perf tool to build CSV that records some relevant hardware counter values
associated with every available AVX-512 instruction form.

## Output format

The output is printed in CSV format.

There are 6 columns:

1. `extension` is a extension this instruction form belongs to
2. `instruction form` is a combination of operands applied to specific opcode
3. `class` is a category this instruction form can reach
4. `level0` shows `core_power.lvl0_turbo_license` hardware counter value
5. `level1` shows `core_power.lvl1_turbo_license` hardware counter value
6. `level2` shows `core_power.lvl2_turbo_license` hardware counter value

Example output line:

```
"avx512f","KANDNW K, K, K","turbo0","1249200","0","0"
```

Example of the complete output is provided in [avx512_core_i9_7900x.csv](/avx512_core_i9_7900x.csv) file.

> **Disclaimer**: provided example is not a reliable reference. The results may vary between
> collector runs, execution on different machines may lead to other results as well.

## Requirements

* Go 1.11 or above (AVX-512 support)
* Linux perf that recognizes `core_power.lvl{0,1,2}_turbo_license` events
* Intel CPU with at least `avx512f`

> Hint: [pmu-tools](https://github.com/andikleen/pmu-tools) contains ocperf.py
> that can be used on systems with older `perf` that does not recognize
> required CPU events even if machine has them.

## Usage

```bash
go get -u github.com/Quasilyte/avx512counters
```

The `$GOPATH/bin` is expected to be included into your system `$PATH`.
If it's not, you may want to move installed binary somewhere where it
will be accessible.

`avx512counters -help` output:

```
  -extensions string
    	comma-separated list of extensions to be evaluated (default "avx512f,avx512dq,avx512cd,avx512bw")
  -iformSpanSize uint
    	how many instruction lines form a single iform span. Higher values slow down the collection (default 100)
  -loopCount uint
    	how many times to execute every iform span. Higher values slow down the collection (default 1000000)
  -perf string
    	perf tool binary name. ocperf and other drop-in replacements will do (default "perf")
  -perfRounds uint
    	how many times to re-validate perf results. Higher values slow down the collection (default 1)
  -workDir string
    	where to put results and the intermediate files (default "./avx512counters-workdir")
```

The only thing you might want to adjust is `extensions` argument.

Suppose you're only interested in `avx512f`, then you can run `avx512counters` like this:

```
avx512counters -extension=avx512f | tee results.csv
```

The result CSV is printed to stdout.
Collection status is printed to stderr.
