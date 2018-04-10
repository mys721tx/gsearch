package main

import (
	"bufio"
	"bytes"
	"github.com/biogo/biogo/seq/linear"
	"testing"
)

func TestReadSeq(t *testing.T) {
	seq := "AAAA"

	f := bytes.NewBufferString(">Foo\n" + seq + "\n")

	s := readSeq(bufio.NewReader(f))

	if s.Seq.String() != seq {
		t.Fatalf("expecting %s, readSeq returns %s", seq, s.Seq.String())
	}

}

func TestScanSeq(t *testing.T) {

	seq1 := "AAAA"
	seq2 := "BBBB"

	f := bytes.NewBufferString(">Foo\n" + seq1 + "\n>Bar\n" + seq2 + "\n")

	cSeqs := make(chan *linear.Seq)

	wg.Add(1)

	go scanSeq(bufio.NewReader(f), cSeqs)

	s := <-cSeqs

	if s.Seq.String() != seq1 {
		t.Fatalf("expecting %s, readSeq sends %s", seq1, s.Seq.String())
	}

	s = <-cSeqs

	if s.Seq.String() != seq2 {
		t.Fatalf("expecting %s, readSeq sends %s", seq2, s.Seq.String())
	}

	wg.Wait()
}
