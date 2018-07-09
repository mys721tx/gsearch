/*
 *  GSEARCH: A concurrent tool suite for metagenomics
 *  Copyright (C) 2018  Yishen Miao
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package seqio_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/seqio"
)

type TestSequence struct {
	name string
	seq  string
}

func joinSeq(seqs ...TestSequence) *bytes.Buffer {
	s := ""

	for _, seq := range seqs {
		s += fmt.Sprintf(">%s\n%s\n", seq.name, seq.seq)
	}

	return bytes.NewBufferString(s)
}

func (t *TestSequence) AssertEqual(test *testing.T, s *linear.Seq) {
	if s.Annotation.ID != t.name {
		test.Errorf(
			"expecting name %s, ReadSeq returns %s", t.name, s.Annotation.ID)
	}

	if s.Seq.String() != t.seq {
		test.Errorf("expecting sequence %s, ReadSeq returns %s", t.seq, s.Seq.String())
	}
}

func newTestSeq(s *linear.Seq) TestSequence {
	return TestSequence{name: s.Annotation.ID, seq: s.Seq.String()}
}

func TestReadSeq(t *testing.T) {
	seq := TestSequence{name: "Foo", seq: "AAAA"}

	f := joinSeq(seq)

	s := seqio.ReadSeq(f)

	seq.AssertEqual(t, s)

}

func TestScanSeq(t *testing.T) {

	var wg sync.WaitGroup

	seq1 := TestSequence{name: "Foo", seq: "AAAA"}
	seq2 := TestSequence{name: "Bar", seq: "GGGG"}

	f := joinSeq(seq1, seq2)

	out := make(chan *linear.Seq)

	wg.Add(1)

	go seqio.ScanSeq(f, out, &wg)

	s := <-out

	seq1.AssertEqual(t, s)

	s = <-out

	seq2.AssertEqual(t, s)

	wg.Wait()
}

func TestWriteSeq(t *testing.T) {

	var wg sync.WaitGroup

	seq1 := linear.NewSeq("Foo", []alphabet.Letter("AAAA"), alphabet.DNA)
	seq2 := linear.NewSeq("Bar", []alphabet.Letter("GGGG"), alphabet.DNA)

	fExpected := joinSeq(newTestSeq(seq1), newTestSeq(seq2))

	in := make(chan *linear.Seq)

	wg.Add(1)

	f := bytes.NewBufferString("")

	go seqio.WriteSeq(f, in, &wg)

	in <- seq1
	in <- seq2

	close(in)

	wg.Wait()

	if f.String() != fExpected.String() {
		t.Errorf("expecting sequence %s, WriteSeq returns %s", fExpected.String(), f.String())
	}
}
