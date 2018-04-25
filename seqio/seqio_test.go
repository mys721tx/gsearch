package seqio

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/biogo/biogo/seq/linear"
	"sync"
	"testing"
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
		test.Errorf("expecting name %s, readSeq returns %s", t.name, s.Annotation.ID)
	}

	if s.Seq.String() != t.seq {
		test.Errorf("expecting sequence %s, readSeq returns %s", t.seq, s.Seq.String())
	}
}

func TestReadSeq(t *testing.T) {
	seq := TestSequence{name: "Foo", seq: "AAAA"}

	f := joinSeq(seq)

	s := ReadSeq(bufio.NewReader(f))

	seq.AssertEqual(t, s)

}

func TestScanSeq(t *testing.T) {

	var wg sync.WaitGroup

	seq1 := TestSequence{name: "Foo", seq: "AAAA"}
	seq2 := TestSequence{name: "Bar", seq: "GGGG"}

	f := joinSeq(seq1, seq2)

	cSeqs := make(chan *linear.Seq)

	wg.Add(1)

	go ScanSeq(bufio.NewReader(f), cSeqs, &wg)

	s := <-cSeqs

	seq1.AssertEqual(t, s)

	s = <-cSeqs

	seq2.AssertEqual(t, s)

	wg.Wait()
}
