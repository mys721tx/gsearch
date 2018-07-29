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

package derep_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"

	"github.com/stretchr/testify/assert"

	"github.com/mys721tx/gsearch/pkg/derep"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

func TestParseAnno(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=100",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 100, "Size should be the value of size.")
}
func TestParseAnnoMissingSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is missing.")
}

func TestParseAnnoUnrecognizedSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=spam",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is not a int.")
}

func TestParseAnnoNegativeSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=-200",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is negative.")
}

func TestParseAnnoMultipleSizes(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=100;size=200",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 200, "Size should be the last value of size.")
}

func TestParseAnnoMissingName(t *testing.T) {
	seq := linear.NewSeq(
		"size=100",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "sequence",
		"Name should default to sequence.",
	)
	assert.Equal(t, res.Size, 100, "Size should be the value of size.")
}

func TestParseAnnoMultipleMonads(t *testing.T) {
	seq := linear.NewSeq(
		"size=100;foo;spam=egg;bar",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 100, "Size should be the value of size.")
}

func TestDeRep(t *testing.T) {

	var wg sync.WaitGroup

	seqs := []*linear.Seq{
		linear.NewSeq(
			"foo;size=100",
			[]alphabet.Letter("ATTC"),
			alphabet.DNA,
		),
		linear.NewSeq(
			"bar;size=10",
			[]alphabet.Letter("ATTC"),
			alphabet.DNA,
		),
		linear.NewSeq(
			"size=4;spam",
			[]alphabet.Letter("ATTC"),
			alphabet.DNA,
		),
	}

	c := make(chan *linear.Seq)

	w := new(bytes.Buffer)

	wg.Add(1)

	go derep.DeRep(c, w, &wg)

	for _, seq := range seqs {
		c <- seq
	}

	close(c)

	wg.Wait()

	seq, _ := seqio.ReadSeq(w)

	res := derep.ParseAnno(seq)

	assert.Equal(t, res.Name, "foo",
		"Name should be the first monad of the first sequence.",
	)
	assert.Equal(t, res.Size, 114, "Size should be the sum of all sizes.")
}
