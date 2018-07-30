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

package seqio_test

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mys721tx/gsearch/mocks"

	"github.com/mys721tx/gsearch/pkg/seqio"
)

var wg sync.WaitGroup

// writeString writes the sequence to string via a writer
func writeString(s *linear.Seq) string {
	f := new(bytes.Buffer)

	w := fasta.NewWriter(f, seqio.WidthCol)

	w.Write(s)

	return f.String()
}

// assertEqualSeq is a wrapper for asserting the equalvance of two linear.Seq
func assertEqualSeq(t *testing.T, a, b *linear.Seq) {
	assert.Equal(t, a.ID, b.ID,
		"ID should be the same as input",
	)
	assert.Equal(t, a.Desc, b.Desc,
		"Desc should be the same as input.",
	)
	assert.Equal(t, a.Loc, b.Loc,
		"Loc should be the same as input",
	)
	assert.Equal(t, a.Strand, b.Strand,
		"Strand should be the same as input.",
	)
	assert.Equal(t, a.Conform, b.Conform,
		"Conform should be the same as input.",
	)
	assert.Equal(t, a.Offset, b.Offset,
		"Offset should be the same as input.",
	)
	assert.Equal(t, a.Seq, b.Seq,
		"Seq should be the FASTA sequence.",
	)
}

func TestReadSeq(t *testing.T) {
	seq := linear.NewSeq("Foo", []alphabet.Letter("AAAA"), alphabet.DNA)

	f := bytes.NewBufferString(writeString(seq))

	s, err := seqio.ReadSeq(f)

	if assert.NoError(t, err) {
		assertEqualSeq(t, s, seq)
	}
}

func TestReadSeqReaderError(t *testing.T) {
	f := new(mocks.Reader)

	errs := []error{
		os.ErrPermission,
		os.ErrNotExist,
		os.ErrClosed,
		io.EOF,
	}

	for _, e := range errs {
		f.On("Read", mock.Anything).Return(0, e)

		s, err := seqio.ReadSeq(f)

		if assert.Error(t, err) {
			assert.Nil(t, s,
				"nil should be returned when an error occurs.",
			)
		}
	}
}

func TestReadSeqMalform(t *testing.T) {
	f := bytes.NewBufferString("AAAA\n")

	s, err := seqio.ReadSeq(f)

	if assert.Error(t, err) {
		assert.Nil(t, s,
			"nil should be returned when an error occurs.",
		)
	}
}

func TestScanSeq(t *testing.T) {
	seqs := []*linear.Seq{
		linear.NewSeq("Foo", []alphabet.Letter("AAAA"), alphabet.DNA),
		linear.NewSeq("Bar", []alphabet.Letter("GGGG"), alphabet.DNA),
	}

	var fExp string

	for _, s := range seqs {
		fExp += writeString(s)
	}

	c := make(chan *linear.Seq)

	f := bytes.NewBufferString(fExp)

	wg.Add(1)

	go seqio.ScanSeq(f, c, &wg)

	for _, seq := range seqs {
		s := <-c

		assertEqualSeq(t, s, seq)
	}

	wg.Wait()
}

func TestScanSeqReaderError(t *testing.T) {
	c := make(chan *linear.Seq)
	f := new(mocks.Reader)

	errs := []error{
		os.ErrPermission,
		os.ErrNotExist,
		os.ErrClosed,
		io.EOF,
	}

	for _, err := range errs {
		f.On("Read", mock.Anything).Return(0, err)

		wg.Add(1)

		go assert.Panics(
			t, func() { seqio.ScanSeq(f, c, &wg) },
			"ScanSeq should panic when encountered an error",
		)

		s := <-c

		assert.Nil(t, s,
			"nil should be returned when an error occurs.",
		)
	}

	wg.Wait()
}

func TestScanSeqMalform(t *testing.T) {
	f := bytes.NewBufferString("AAAA\n>Foo\nAAAA")

	c := make(chan *linear.Seq)

	wg.Add(1)

	go assert.Panics(t, func() { seqio.ScanSeq(f, c, &wg) },
		"ScanSeq should panic when encountered an error",
	)

	s := <-c

	assert.Nil(t, s,
		"nil should be returned when an error occurs.",
	)

	wg.Wait()
}

func TestWriteSeq(t *testing.T) {
	seqs := []*linear.Seq{
		linear.NewSeq("Foo", []alphabet.Letter("AAAA"), alphabet.DNA),
		linear.NewSeq("Bar", []alphabet.Letter("GGGG"), alphabet.DNA),
	}

	c := make(chan *linear.Seq)

	wg.Add(1)

	f := new(bytes.Buffer)

	go seqio.WriteSeq(f, c, &wg)

	var fExp string

	for _, s := range seqs {
		c <- s
		fExp += writeString(s)
	}

	close(c)

	wg.Wait()

	assert.Equal(t, f.String(), fExp,
		"Output should be the same as input sequence.",
	)
}

func TestWriteSeqWriterError(t *testing.T) {
	f := new(mocks.Writer)
	seq := linear.NewSeq("Foo", []alphabet.Letter("AAAA"), alphabet.DNA)

	errs := []error{
		os.ErrPermission,
		os.ErrNotExist,
		os.ErrClosed,
		io.EOF,
	}

	for _, err := range errs {
		f.On("Write", mock.Anything).Return(0, err)

		wg.Add(1)

		c := make(chan *linear.Seq)

		go assert.Panics(
			t, func() { seqio.WriteSeq(f, c, &wg) },
			"WriteSeq should panic when its writer encounters an error.",
		)

		c <- seq

		close(c)
	}

	wg.Wait()
}
