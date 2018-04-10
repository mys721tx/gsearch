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
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"io"
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

func readSeq(f io.Reader) *linear.Seq {
	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	s, err := r.Read()

	if err != nil {
		log.Fatalf("failed to read sequence: %v", err)
	}

	// Type assertion to linear.Seq
	return s.(*linear.Seq)
}

func scanSeq(f io.Reader, seqs chan<- *linear.Seq) {
	defer wg.Done()

	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	sc := seqio.NewScanner(r)

	for sc.Next() {
		s := sc.Seq()
		err := sc.Error()

		if err != nil {
			log.Fatalf("failed during read: %v", err)
		} else {
			// Type assertion to linear.Seq
			seqs <- s.(*linear.Seq)
		}
	}
	close(seqs)
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

	smith := align.SWAffine{
		Matrix:  *makeScoreMatrix(),
		GapOpen: *gapopen,
	}

	fRef, err := os.Open(*ref)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *ref, err)
	}

	sRef := readSeq(fRef)

	fTgt, err := os.Open(*tgt)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *tgt, err)
	}

	csTgt := make(chan *linear.Seq)

	wg.Add(2)
	go scanSeq(fTgt, csTgt)
	go alignSW(sRef, smith, csTgt)

	wg.Wait()
}
