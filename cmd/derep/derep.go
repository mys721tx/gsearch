// GSEARCH: A concurrent tool suite for metagenomics
// Copyright (C) 2018  Yishen Miao
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/cluster"
	"github.com/mys721tx/gsearch/pkg/derep"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

var (
	pin, pout string
	max, min  int
	wg        sync.WaitGroup
)

func main() {

	flag.IntVar(
		&min,
		"min",
		cluster.MinLen,
		"minimal abundance of a sequence, default to 0.",
	)

	flag.IntVar(
		&max,
		"max",
		cluster.MaxLen,
		"maximal abundance of a sequence, default to 0.",
	)

	flag.StringVar(
		&pin,
		"in",
		"",
		"path to the sequence FASTA file, default to stdin.",
	)

	flag.StringVar(
		&pout,
		"out",
		"",
		"path to the output FASTA file, default to stdout.",
	)
	flag.Parse()

	var fin, fout *os.File

	if pin == "" {
		fin = os.Stdin
	} else if f, err := os.Open(pin); err == nil {
		fin = f
	} else {
		log.Panicf("failed to open %q: %v", pin, err)
	}

	if pout == "" {
		fout = os.Stdout
	} else if f, err := os.Create(pout); err == nil {
		fout = f
	} else {
		log.Panicf("failed to open %q: %v", pout, err)
	}

	defer func() {
		if err := fout.Close(); err != nil {
			log.Panicf("failed to close %q: %v", pout, err)
		}
	}()

	w := bufio.NewWriter(fout)

	defer func() {
		if err := w.Flush(); err != nil {
			log.Panicf("failed to flush %q: %v", pout, err)
		}
	}()

	c := make(chan *linear.Seq)

	wg.Add(2)
	go seqio.ScanSeq(fin, c, &wg)       // TODO: handling panic
	go derep.DeRep(c, w, min, max, &wg) // TODO: handling panic
	wg.Wait()
}
