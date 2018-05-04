/*
 *
 * Copyright 2018 Yishen Miao
 *
 */

package seqio

import (
	"io"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
	"github.com/golang/glog"
)

// ReadSeq reads a sequence from a fasta file.
func ReadSeq(f io.Reader) *linear.Seq {
	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	s, err := r.Read()

	if err != nil {
		glog.Fatalf("failed to read sequence: %v", err)
	}

	// Type assertion to linear.Seq
	seq := s.(*linear.Seq)

	if v, p := seq.Validate(); !v {
		glog.Fatalf("invalidate symbol: position %v", p)
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

		if err := sc.Error(); err != nil {
			glog.Warningf("failed during read: %v", err)
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
		if _, err := w.Write(seq); err != nil {
			glog.Warningf("failed to write: %v", err)
		}
	}
}
