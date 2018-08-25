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

package cluster_test

import (
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"

	"github.com/stretchr/testify/assert"

	"github.com/mys721tx/gsearch/pkg/cluster"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

func parseBuf(r io.Reader) *cluster.Cluster {
	seq, _ := seqio.ReadSeq(r)

	return cluster.ParseAnno(seq)
}

func TestClusterName(t *testing.T) {
	c := cluster.Cluster{
		Seq: *linear.NewSeq(
			"foo",
			[]alphabet.Letter("AAAA"),
			alphabet.DNA,
		),
		Size: 100,
	}

	res := c.Name()

	assert.Equal(
		t, fmt.Sprintf("%v;size=%d", c.ID, c.Size), res,
		"Name should return a semicolon seperated list of ID and size.",
	)
}

func TestParseAnno(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=100",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 100, "Size should be the value of size.")
}
func TestParseAnnoMissingSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is missing.")
}

func TestParseAnnoUnrecognizedSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=spam",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is not a int.")
}

func TestParseAnnoNegativeSize(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=-200",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 1, "Size should be 1 when size is negative.")
}

func TestParseAnnoMultipleSizes(t *testing.T) {
	seq := linear.NewSeq(
		"foo;size=100;size=200",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 200, "Size should be the last value of size.")
}

func TestParseAnnoMissingName(t *testing.T) {
	seq := linear.NewSeq(
		"size=100",
		[]alphabet.Letter("ATTC"),
		alphabet.DNA,
	)

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "sequence",
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

	res := cluster.ParseAnno(seq)

	assert.Equal(t, res.ID, "foo", "Name should be the first monad.")
	assert.Equal(t, res.Size, 100, "Size should be the value of size.")
}

func BenchmarkParseAnno(b *testing.B) {
	seq := linear.NewSeq(
		"size=100;foo;spam=egg;bar",
		[]alphabet.Letter("AAAA"),
		alphabet.DNA,
	)

	for n := 0; n < b.N; n++ {
		cluster.ParseAnno(seq)
	}
}

func TestClusterPassFilter(t *testing.T) {
	seq := cluster.Cluster{
		Size: 100,
	}

	res := seq.PassFilter(10, 1000)

	assert.True(t, res,
		"Cluster that is within the given filter sizes should pass.",
	)
}

func TestClusterMinFilter(t *testing.T) {
	seq := cluster.Cluster{
		Size: 10,
	}

	res := seq.PassFilter(10, 1000)

	assert.True(t, res,
		"Cluster that is strictly smaller the given size should not pass.",
	)

	seq.Size = 9

	res = seq.PassFilter(10, 1000)

	assert.False(t, res,
		"Cluster that is strictly smaller the given size should not pass.",
	)
}

func TestClusterMaxFilter(t *testing.T) {
	seq := cluster.Cluster{
		Size: 1000,
	}

	res := seq.PassFilter(10, 1000)

	assert.True(t, res,
		"Cluster that is strictly greater the given size should not pass.",
	)

	seq.Size = 1001

	res = seq.PassFilter(10, 1000)

	assert.False(t, res,
		"Cluster that is strictly greater the given size should not pass.",
	)
}

func TestByAbundanceSize(t *testing.T) {
	clus := []*cluster.Cluster{
		{
			Seq: *linear.NewSeq(
				"foo",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 100,
		},
		{
			Seq: *linear.NewSeq(
				"bar",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 10,
		},
		{
			Seq: *linear.NewSeq(
				"spam",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 1,
		},
		{
			Seq: *linear.NewSeq(
				"egg",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 1000,
		},
	}

	expects := []int{1000, 100, 10, 1}

	sort.Sort(cluster.ByAbundance(clus))

	for i, v := range expects {
		assert.Equal(t, v, clus[i].Size,
			"Cluster are sorted by their abundance from high to low.",
		)
	}
}

func TestByAbundanceName(t *testing.T) {
	clus := []*cluster.Cluster{
		{
			Seq: *linear.NewSeq(
				"B",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 100,
		},
		{
			Seq: *linear.NewSeq(
				"D",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 100,
		},
		{
			Seq: *linear.NewSeq(
				"A",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 100,
		},
		{
			Seq: *linear.NewSeq(
				"C",
				[]alphabet.Letter("AAAA"),
				alphabet.DNA,
			),
			Size: 100,
		},
	}

	expects := []string{"A", "B", "C", "D"}

	sort.Sort(cluster.ByAbundance(clus))

	for i, v := range expects {
		assert.Equal(t, v, clus[i].ID,
			"Clusters of the same size are sorted by name in lexicographical order.",
		)
	}
}
