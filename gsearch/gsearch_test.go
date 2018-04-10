package main

import (
	"bufio"
	"bytes"
	"github.com/biogo/biogo/seq/linear"
	"testing"
)

type TestSequence struct {
	name string
	seq  string
}

func (t *TestSequence) AssertEqual(test *testing.T, s *linear.Seq) {
	if s.Annotation.ID != t.name {
		test.Errorf("expecting name %s, readSeq returns %s", t.name, s.Annotation.ID)
	}

	if s.Seq.String() != t.seq {
		test.Errorf("expecting sequence %s, readSeq returns %s", t.seq, s.Seq.String())
	}
}

func TestReadSeq(t *testing.T) {
	seq := TestSequence{name: "Foo", seq: "AAAA"}

	f := bytes.NewBufferString(">" + seq.name + "\n" + seq.seq + "\n")

	s := readSeq(bufio.NewReader(f))

	seq.AssertEqual(t, s)

}

func TestScanSeq(t *testing.T) {

	seq1 := TestSequence{name: "Foo", seq: "AAAA"}
	seq2 := TestSequence{name: "Bar", seq: "GGGG"}

	f := bytes.NewBufferString(">" + seq1.name + "\n" + seq1.seq + "\n>" + seq2.name + "\n" + seq2.seq + "\n")

	cSeqs := make(chan *linear.Seq)

	wg.Add(1)

	go scanSeq(bufio.NewReader(f), cSeqs)

	s := <-cSeqs

	seq1.AssertEqual(t, s)

	s = <-cSeqs

	seq2.AssertEqual(t, s)

	wg.Wait()
}
