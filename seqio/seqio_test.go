package seqio

import (
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
		test.Errorf(
			"expecting name %s, ReadSeq returns %s", t.name, s.Annotation.ID)
	}

	if s.Seq.String() != t.seq {
		test.Errorf("expecting sequence %s, ReadSeq returns %s", t.seq, s.Seq.String())
	}
}

func newSeqLinear(s *linear.Seq) TestSequence {
	return TestSequence{name: s.Annotation.ID, seq: s.Seq.String()}
}

func TestReadSeq(t *testing.T) {
	seq := TestSequence{name: "Foo", seq: "AAAA"}

	f := joinSeq(seq)

	s := ReadSeq(f)

	seq.AssertEqual(t, s)

}

func TestScanSeq(t *testing.T) {

	var wg sync.WaitGroup

	seq1 := TestSequence{name: "Foo", seq: "AAAA"}
	seq2 := TestSequence{name: "Bar", seq: "GGGG"}

	f := joinSeq(seq1, seq2)

	out := make(chan *linear.Seq)

	wg.Add(1)

	go ScanSeq(f, out, &wg)

	s := <-out

	seq1.AssertEqual(t, s)

	s = <-out

	seq2.AssertEqual(t, s)

	wg.Wait()
}
