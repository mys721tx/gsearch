# GSEARCH

[![GoDoc](https://godoc.org/github.com/mys721tx/gsearch?status.svg)](https://godoc.org/github.com/mys721tx/gsearch)

[vsearch](https://github.com/torognes/vsearch) implemented in Go.

## Description

**GSEARCH** is a concurrent tool suite for metagenomics. Currently only the
derepuliator is implemented.

## Installation

1. Clone this repository: `git clone https://github.com/mys721tx/gsearch.git`
2. Run `dep ensure` to install dependencies.
    * See [Install Instruction](https://golang.org/doc/install) to install Go.
        GSEARCH has been tested on Go 1.10.2.
    * See [Installation](https://golang.github.io/dep/docs/installation.html) to
        install `dep`.
3. Run `go build cmd/derep/derep.go` to build the binary for `derep`.
4. Run `./derep -in infile -out outfile` to remove duplicated FASTA sequence
    from `infile` and output to `outfile`.
    * The description for a FASTA sequence should be in `NAME;size=NUM` format
        where `NAME` is a string and `NUM` is a integer.

## Testing Dataset

1. Follow the [VSEARCH pipeline](https://github.com/torognes/vsearch/wiki/VSEARCH-pipeline)
    demo until `all.fasta` is generated.
    * See [README](https://github.com/torognes/vsearch) to install VSEARCH.
2. Run `./derep -in all.fasta -out all.derep.fasta` to generate GSEARCH output.
3. Follow the pipeline to generate `all.derep.fasta` from VSEARCH.
4. Run `grep ">" all.derep.fasta | sort > sorted.fasta` to generate sorted
    descriptions for both outputs.
5. Compare the output using a diff tool. Compare the running time and memory
    usage using a profiler.

## Author

* [Yishen Miao](https://github.com/mys721tx)

## License

[GNU General Public License, version 3](http://www.gnu.org/licenses/gpl-3.0.html)
