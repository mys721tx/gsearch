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
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"log"
	"os"
)

var (
	ref      = flag.String("reference", "", "reference sequence to align")
	tgt      = flag.String("target", "", "target sequence to align")
	match    = flag.Int("match", 2, "score for match")
	mismatch = flag.Int("mismatch", -1, "score for mismatch")
	gap      = flag.Int("gap", -2, "score for gap")
	gapopen  = flag.Int("gap_open", 0, "score for gap open")
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

func readSeq(f *os.File) *linear.Seq {

	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	s, err := r.Read()

	if err != nil {
		log.Fatalf("failed to read sequence: %v", err)
	}

	// Type assertion to linear.Seq
	return s.(*linear.Seq)
}

func main() {
	flag.Parse()

	fRef, err := os.Open(*ref)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *ref, err)
	}

	sRef := readSeq(fRef)

	fTgt, err := os.Open(*tgt)

	if err != nil {
		log.Fatalf("failed to open %q: %v", *ref, err)
	}

	sTgt := readSeq(fTgt)

	smith := align.SWAffine{
		Matrix:  *makeScoreMatrix(),
		GapOpen: *gapopen,
	}

	aln, err := smith.Align(sRef, sTgt)

	if err != nil {
		log.Fatalf("failed to align: %v", err)
	}

	fmt.Printf("%s\n", aln)
	fa := align.Format(sRef, sTgt, aln, '-')
	fmt.Printf("%s\n%s\n", fa[0], fa[1])
}
