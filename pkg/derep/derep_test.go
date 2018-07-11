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
package derep_test

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/derep"
)

// TestSeq use mnd for monadic items, par for key-value pairs, and seq for
// sequence.
type TestSeq struct {
	mnd []string
	par map[string]string
	seq string
}

func (t *TestSeq) newSeqLinear() *linear.Seq {
	items := t.mnd
	for k, v := range t.par {
		items = append(items, fmt.Sprintf("%s=%s", k, v))
	}

	return linear.NewSeq(
		strings.Join(items, ";"),
		[]alphabet.Letter(t.seq),
		alphabet.DNA,
	)
}

func TestParseAnno(t *testing.T) {
	exp := derep.Cluster{
		Name: "foo",
		Size: 100,
	}

	seq := TestSeq{
		mnd: []string{exp.Name},
		par: map[string]string{"size": "100"},
		seq: "ATTC",
	}

	res := derep.ParseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}
func TestParseAnnoMissingSize(t *testing.T) {
	exp := derep.Cluster{
		Name: "foo",
		Size: 1,
	}

	seq := TestSeq{
		mnd: []string{exp.Name},
		par: map[string]string{},
		seq: "ATTC",
	}

	res := derep.ParseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoUnrecognizedSize(t *testing.T) {
	exp := derep.Cluster{
		Name: "foo",
		Size: 1,
	}

	seq := TestSeq{
		mnd: []string{exp.Name},
		par: map[string]string{"size": "spam"},
		seq: "ATTC",
	}

	res := derep.ParseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoNegativeSize(t *testing.T) {
	exp := derep.Cluster{
		Name: "foo",
		Size: 1,
	}

	seq := TestSeq{
		mnd: []string{exp.Name},
		par: map[string]string{"size": "-200"},
		seq: "ATTC",
	}

	res := derep.ParseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoMissingName(t *testing.T) {
	exp := derep.Cluster{
		Name: "sequence",
		Size: 100,
	}

	seq := TestSeq{
		mnd: []string{},
		par: map[string]string{"size": "100"},
		seq: "ATTC",
	}

	res := derep.ParseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestDeRep(t *testing.T) {

	var wg sync.WaitGroup

	exp := derep.Cluster{
		Name: "foo",
		Size: 114,
	}

	seqs := []TestSeq{
		{
			mnd: []string{"foo"},
			par: map[string]string{"size": "100"},
			seq: "ATTC",
		},
		{
			mnd: []string{"bar"},
			par: map[string]string{"size": "10"},
			seq: "ATTC",
		},
		{
			mnd: []string{},
			par: map[string]string{"size": "4"},
			seq: "ATTC",
		},
	}

	cin := make(chan *linear.Seq)
	cout := make(chan *linear.Seq)

	wg.Add(1)
	go derep.DeRep(cin, cout, &wg)

	for _, seq := range seqs {
		cin <- seq.newSeqLinear()
	}

	close(cin)

	res := derep.ParseAnno(<-cout)

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}

	wg.Wait()
}
