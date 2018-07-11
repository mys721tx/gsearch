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
	pin = flag.String(
		"in",
		"",
		"path to the sequence fasta file, default to stdin.",
	)
	pout = flag.String(
		"out",
		"",
		"path to the output fasta file, default to stdout.",
	)
	wg sync.WaitGroup
)

type cluster struct {
	name string
	size int
	//seqs *linear.Seq
}

// parseAnno parses the annotation in sequence. The first monad is used as the
// name of the sequence; otherwise defaults to "sequence". The last key-value
// pair with key "size" is used as the size; otherwise defaults to 1.
func parseAnno(seq *linear.Seq) *cluster {

	var monads []string

	pairs := make(map[string]string)

	for _, item := range strings.Split(seq.Annotation.ID, ";") {
		// If more than 2 then skip
		switch pair := strings.Split(item, "="); len(pair) {
		case 1:
			{
				monads = append(monads, pair[0])
			}
		case 2:
			{
				pairs[pair[0]] = pair[1]
			}
		}
	}

	var size int
	var err error

	if _, prs := pairs["size"]; prs {
		size, err = strconv.Atoi(pairs["size"])
		if err != nil {
			size = 1
		}
		if size <= 0 {
			size = 1
		}
	} else {
		size = 1
	}

	var name string

	if len(monads) > 0 {
		name = monads[0]
	} else {
		name = "sequence"
	}

	res := cluster{name: name, size: size}

	return &res
}

// derep receives a sequence from a channel and builds a map.
// If a sequence is in a map, if not, output it.
func deRep(in <-chan *linear.Seq, out chan<- *linear.Seq) {
	defer wg.Done()
	defer close(out)

	rep := make(map[string]*cluster)

	for seq := range in {
		c := parseAnno(seq)

		if _, prs := rep[seq.String()]; !prs {
			//rep[s] = &cluster{name: name, size: size, seqs: seq}
			rep[seq.String()] = c
		} else {
			rep[seq.String()].size += c.size
		}
	}

	for s, v := range rep {
		out <- linear.NewSeq(
			fmt.Sprintf("%v;size=%d", v.name, v.size),
			[]alphabet.Letter(s),
			alphabet.DNA,
		)
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

	defer func() {
		err = fout.Close()

		if err != nil {
			log.Fatalf("failed to close %q: %v", *pout, err)
		}
	}()

	w := bufio.NewWriter(fout)

	cin := make(chan *linear.Seq)
	cout := make(chan *linear.Seq)

	wg.Add(3)
	go seqio.ScanSeq(fin, cin, &wg)
	go deRep(cin, cout)
	go seqio.WriteSeq(w, cout, &wg)
	wg.Wait()

	err = w.Flush()

	if err != nil {
		log.Fatalf("failed to flush %q: %v", *pout, err)
	}

}
