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
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/seqio"
)

var (
	ref      = flag.String("reference", "", "path to the reference sequence fasta file")
	tgt      = flag.String("target", "", "path to the target sequence fasta file")
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
