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
	"github.com/biogo/biogo/seq/linear"
	"log"
)

var (
	reference = flag.String("reference", "", "reference sequence to align")
	target    = flag.String("target", "", "target sequence to align")
	match     = flag.Int("match", 2, "score for match")
	mismatch  = flag.Int("mismatch", -1, "score for mismatch")
	gap       = flag.Int("gap", -2, "score for gap")
	gapopen   = flag.Int("gap_open", 0, "score for gap open")
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

func main() {
	flag.Parse()

	seqReference := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(*reference))}
	seqReference.Alpha = alphabet.DNAgapped
	seqTarget := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(*target))}
	seqTarget.Alpha = alphabet.DNAgapped

	smith := align.SWAffine{
		Matrix:  *makeScoreMatrix(),
		GapOpen: *gapopen,
	}

	aln, err := smith.Align(seqReference, seqTarget)

	if err != nil {
		log.Fatalf("failed to align: %v", err)
	}

	fmt.Printf("%s\n", aln)
	fa := align.Format(seqReference, seqTarget, aln, '-')
	fmt.Printf("%s\n%s\n", fa[0], fa[1])
}
