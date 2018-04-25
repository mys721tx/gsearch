/*
 *
 * Copyright 2018 Yishen Miao
 *
 */

package seqio

import (
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"io"
	"log"
	"sync"
)

// ReadSeq reads a sequence from a fasta file.
func ReadSeq(f io.Reader) *linear.Seq {
	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	s, err := r.Read()

	if err != nil {
		log.Fatalf("failed to read sequence: %v", err)
	}

	// Type assertion to linear.Seq
	seq := s.(*linear.Seq)

	if v, p := seq.Validate(); !v {
		log.Fatalf("invalidate symbol: position %v", p)
	}

	return seq
}

// ScanSeq scans sequences from a fasta file to a channel.
func ScanSeq(f io.Reader, out chan<- *linear.Seq, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	sc := seqio.NewScanner(r)

	for sc.Next() {
		s := sc.Seq()
		err := sc.Error()

		if err != nil {
			log.Fatalf("failed during read: %v", err)
		} else {
			// Type assertion to linear.Seq
			out <- s.(*linear.Seq)
		}
	}
}

// WriteSeq writes sequences from a channel to a fasta file.
func WriteSeq(f io.Writer, in <-chan *linear.Seq, wg *sync.WaitGroup) {
	defer wg.Done()

	w := fasta.NewWriter(f, 80)

	for seq := range in {
		_, err := w.Write(seq)
		if err != nil {
			log.Fatalf("failed to write: %v", err)
		}
	}
}
