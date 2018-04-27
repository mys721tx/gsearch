/*
 *
 * Copyright 2018 Yishen Miao
 *
 */

package main

import (
	"bufio"
	"flag"
	"github.com/biogo/biogo/seq/linear"
	"github.com/mys721tx/gsearch/seqio"
	"log"
	"os"
	"sync"
)

var (
	pin  = flag.String("in", "", "path to the sequence fasta file")
	pout = flag.String("out", "", "path to the output fasta file")
	wg   sync.WaitGroup
)

// derep receives a sequence from a channel and builds a map.
// If a sequence is in a map, if not, output it.
func deRep(in <-chan *linear.Seq, out chan<- *linear.Seq) {
	defer wg.Done()
	defer close(out)

	rep := make(map[string]bool)

	for seq := range in {
		if !rep[seq.Seq.String()] {
			rep[seq.Seq.String()] = true
			out <- seq
		}
	}
}

func main() {
	flag.Parse()

	fin, err := os.Open(*pin)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *pin, err)
	}

	fout, err := os.Create(*pout)

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
