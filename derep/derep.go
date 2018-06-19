/*
 *  GSEARCH: A concurrent tool suite for metagenomics
 *  Copyright (C) 2018  Yishen Miao
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
	"github.com/mys721tx/gsearch/seqio"
)

var (
	pin  = flag.String("in", "", "path to the sequence fasta file")
	pout = flag.String("out", "", "path to the output fasta file")
	wg   sync.WaitGroup
)

type cluster struct {
	name string
	size int
	//seqs *linear.Seq
}

func parseSeq(seq *linear.Seq) (string, int, string) {
	// hard coded for now
	desc := strings.Split(seq.Annotation.ID, ";size=")
	i, _ := strconv.ParseInt(desc[1], 10, 0)
	s := seq.String()
	return desc[0], int(i), s
}

// derep receives a sequence from a channel and builds a map.
// If a sequence is in a map, if not, output it.
func deRep(in <-chan *linear.Seq, out chan<- *linear.Seq) {
	defer wg.Done()
	defer close(out)

	rep := make(map[string]*cluster)

	for seq := range in {
		name, size, s := parseSeq(seq)

		if _, prs := rep[s]; !prs {
			//rep[s] = &cluster{name: name, size: size, seqs: seq}
			rep[s] = &cluster{name: name, size: size}
		} else {
			rep[s].size += size
		}
	}

	for s, v := range rep {
		out <- linear.NewSeq(
			fmt.Sprintf("%v;size=%d", v.name, v.size),
			[]alphabet.Letter(s),
			alphabet.DNA)
	}

}

func main() {
	flag.Parse()

	var fin, fout *os.File
	var err error

	if *pin == "" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(*pin)
	}

	if err != nil {
		log.Fatalf("failed to open %q: %v", *pin, err)
	}

	if *pout == "" {
		fout = os.Stdin
	} else {
		fout, err = os.Create(*pout)
	}

	if err != nil {
		log.Fatalf("failed to open %q: %v", *pout, err)
	}

	defer fout.Close()

	w := bufio.NewWriter(fout)
	defer w.Flush()

	cin := make(chan *linear.Seq)
	cout := make(chan *linear.Seq)

	wg.Add(3)
	go seqio.ScanSeq(fin, cin, &wg)
	go deRep(cin, cout)
	go seqio.WriteSeq(w, cout, &wg)
	wg.Wait()
}
