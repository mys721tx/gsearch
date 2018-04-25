/*
 *
 * Copyright 2018 Yishen Miao
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/seq/linear"
	"github.com/mys721tx/gsearch/seqio"
	"log"
	"os"
	"sync"
)

var (
	ref      = flag.String("reference", "", "path to the reference sequence fasta file")
	tgt      = flag.String("target", "", "path to the target sequence fasta file")
	output   = flag.String("output", "", "path to the output fasta file")
	match    = flag.Int("match", 2, "score for match")
	mismatch = flag.Int("mismatch", -1, "score for mismatch")
	gap      = flag.Int("gap", -2, "score for gap")
	gapopen  = flag.Int("gap_open", 0, "score for gap open")
	wg       sync.WaitGroup
)

func makeScoreMatrix() *align.Linear {
	m := align.Linear{
		{0, *gap, *gap, *gap, *gap},
		{*gap, *match, *mismatch, *mismatch, *mismatch},
		{*gap, *mismatch, *match, *mismatch, *mismatch},
		{*gap, *mismatch, *mismatch, *match, *mismatch},
		{*gap, *mismatch, *mismatch, *mismatch, *match},
	}
	return &m
}

func alignSW(ref *linear.Seq, score align.NWAffine, seqs <-chan *linear.Seq) {
	defer wg.Done()

	for tgt := range seqs {
		aln, err := score.Align(ref, tgt)

		if err != nil {
			log.Fatalf("failed to align: %v", err)
		}

		fmt.Printf("%s\n", aln)
		fa := align.Format(ref, tgt, aln, '-')
		fmt.Printf("%s\n%s\n", fa[0], fa[1])
	}
}

func main() {
	flag.Parse()

	nw := align.NWAffine{
		Matrix:  *makeScoreMatrix(),
		GapOpen: *gapopen,
	}

	fRef, err := os.Open(*ref)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *ref, err)
	}

	sRef := seqio.ReadSeq(fRef)

	fTgt, err := os.Open(*tgt)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *tgt, err)
	}

	csTgt := make(chan *linear.Seq)

	wg.Add(2)
	go seqio.ScanSeq(fTgt, csTgt, &wg)
	go alignSW(sRef, nw, csTgt)

	wg.Wait()
}
