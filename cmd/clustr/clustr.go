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
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/derep"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

var (
	pin = flag.String(
		"in",
		"",
		"path to the sequence FASTA file, default to stdin.",
	)
	pout = flag.String(
		"out",
		"",
		"path to the output FASTA file, default to stdout.",
	)
	min = flag.Int(
		"min",
		derep.MinLen,
		"minimal abundance of a sequence, default to 0.",
	)
	max = flag.Int(
		"max",
		derep.MaxLen,
		"maximal abundance of a sequence, default to 0.",
	)
	wg sync.WaitGroup
)

// ByAbundance implements methods to sort a slice of a cluster
type ByAbundance []*derep.Cluster

// Len returns the length of a ByAbundance.
func (c ByAbundance) Len() int { return len(c) }

// Swap swaps two elements in a ByAbundance.
func (c ByAbundance) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less establishes the order between two clusters when sort ByAbundance.
//
// The higher abundance cluster is in front of the lower abundance sequence.
// When two clusters have same abundance, they are sorted by the lexicographical
//order of the labels.
func (c ByAbundance) Less(i, j int) bool {
	if c[i].Size > c[j].Size {
		return true
	} else if c[i].Size < c[j].Size {
		return false
	} else if o := strings.Compare(c[i].ID, c[j].ID); o == -1 {
		return true
	}
	return false
}

func main() {
	flag.Parse()

	var fin, fout *os.File

	if *pin == "" {
		fin = os.Stdin
	} else if f, err := os.Open(*pin); err == nil {
		fin = f
	} else {
		log.Panicf("failed to open %q: %v", *pin, err)
	}

	if *pout == "" {
		fout = os.Stdout
	} else if f, err := os.Create(*pout); err == nil {
		fout = f
	} else {
		log.Panicf("failed to open %q: %v", *pout, err)
	}

	defer func() {
		if err := fout.Close(); err != nil {
			log.Panicf("failed to close %q: %v", *pout, err)
		}
	}()

	w := bufio.NewWriter(fout)

	defer func() {
		if err := w.Flush(); err != nil {
			log.Panicf("failed to flush %q: %v", *pout, err)
		}
	}()

	ch := make(chan *linear.Seq)

	wg.Add(1)
	go seqio.ScanSeq(fin, ch, &wg) // TODO: handling panic

	l := func() []*derep.Cluster {

		var l []*derep.Cluster
		for s := range ch {
			c := derep.ParseAnno(s)
			l = append(l, c)
		}

		sort.Sort(ByAbundance(l))

		return l
	}()

	wg.Wait()

	func(f io.Writer, min, max int) {
		w := fasta.NewWriter(f, seqio.WidthCol)

		for _, c := range l {
			if c.PassFilter(min, max) {
				if _, err := w.Write(c); err != nil {
					log.Panicf("Error occurred during write: %s", err)
				}
			}
		}
	}(w, *min, *max)
}
