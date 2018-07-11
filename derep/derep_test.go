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
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

type TestSequence struct {
	monads []string
	pairs  map[string]string
	seq    string
}

func (t *TestSequence) newSeqLinear() *linear.Seq {
	items := t.monads
	for k, v := range t.pairs {
		items = append(items, fmt.Sprintf("%s=%s", k, v))
	}

	return linear.NewSeq(
		strings.Join(items, ";"),
		[]alphabet.Letter(t.seq),
		alphabet.DNA,
	)
}

func TestParseAnno(t *testing.T) {
	exp := cluster{
		name: "foo",
		size: 100,
	}

	seq := TestSequence{
		monads: []string{exp.name},
		pairs:  map[string]string{"size": "100"},
		seq:    "ATTC",
	}

	res := parseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}
func TestParseAnnoMissingSize(t *testing.T) {
	exp := cluster{
		name: "foo",
		size: 1,
	}

	seq := TestSequence{
		monads: []string{exp.name},
		pairs:  map[string]string{},
		seq:    "ATTC",
	}

	res := parseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoUnrecognizedSize(t *testing.T) {
	exp := cluster{
		name: "foo",
		size: 1,
	}

	seq := TestSequence{
		monads: []string{exp.name},
		pairs:  map[string]string{"size": "spam"},
		seq:    "ATTC",
	}

	res := parseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoNegativeSize(t *testing.T) {
	exp := cluster{
		name: "foo",
		size: 1,
	}

	seq := TestSequence{
		monads: []string{exp.name},
		pairs:  map[string]string{"size": "-200"},
		seq:    "ATTC",
	}

	res := parseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestParseAnnoMissingName(t *testing.T) {
	exp := cluster{
		name: "sequence",
		size: 100,
	}

	seq := TestSequence{
		monads: []string{},
		pairs:  map[string]string{"size": "100"},
		seq:    "ATTC",
	}

	res := parseAnno(seq.newSeqLinear())

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}
}

func TestDeRep(t *testing.T) {
	exp := cluster{
		name: "foo",
		size: 114,
	}

	seqs := []TestSequence{
		{
			monads: []string{"foo"},
			pairs:  map[string]string{"size": "100"},
			seq:    "ATTC",
		},
		{
			monads: []string{"bar"},
			pairs:  map[string]string{"size": "10"},
			seq:    "ATTC",
		},
		{
			monads: []string{},
			pairs:  map[string]string{"size": "4"},
			seq:    "ATTC",
		},
	}

	cin := make(chan *linear.Seq)
	cout := make(chan *linear.Seq)
	wg.Add(1)
	go deRep(cin, cout)

	sizeExpected := 0

	for _, seq := range seqs {
		size, _ := strconv.Atoi(seq.pairs["size"])
		sizeExpected += size
		cin <- seq.newSeqLinear()
	}

	close(cin)

	res := parseAnno(<-cout)

	if !reflect.DeepEqual(exp, *res) {
		t.Errorf(
			"expecting cluster %v, parseAnno returns %v",
			exp,
			res,
		)
	}

	wg.Wait()
}
