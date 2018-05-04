/*
 *
 * Copyright 2018 Yishen Miao
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
