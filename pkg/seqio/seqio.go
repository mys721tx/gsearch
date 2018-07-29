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

// Package seqio provides wrapper functions of biogo.
package seqio

import (
	"io"
	"log"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"
)

const (
	// WidthCol is the length of the column used by WriteSeq.
	WidthCol = 80
)

// ReadSeq reads a sequence from a FASTA file.
//
// If the underlaying reader has encountered any error, ReadSeq will return the
// error.
func ReadSeq(f io.Reader) (*linear.Seq, error) {
	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	s, err := r.Read()

	if err != nil {
		return nil, err
	}

	// Type assertion to linear.Seq
	seq := s.(*linear.Seq)

	return seq, nil
}

// ScanSeq scans sequences from a fasta file to a channel.
//
// If the underlaying reader has encountered any error, ScanSeq will panic as
// the reader can no longer be read.
func ScanSeq(f io.Reader, out chan<- *linear.Seq, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	r := fasta.NewReader(f, linear.NewSeq("", nil, alphabet.DNAgapped))

	sc := seqio.NewScanner(r)

	for sc.Next() {
		s := sc.Seq()
		// Type assertion to linear.Seq
		out <- s.(*linear.Seq)
	}

	if err := sc.Error(); err != nil {
		log.Panicf("Error occurred during scan: %s", err)
	}
}

// WriteSeq writes sequences from a channel to a fasta file.
//
// If the underlaying writer has encountered any error, WriteSeq will panic as
// the writer can no longer be written.
func WriteSeq(f io.Writer, in <-chan *linear.Seq, wg *sync.WaitGroup) {
	defer wg.Done()

	w := fasta.NewWriter(f, WidthCol)

	for seq := range in {
		if _, err := w.Write(seq); err != nil {
			log.Panicf("Error occurred during write: %s", err)
		}
	}
}
