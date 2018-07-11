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
	seq := TestSequence{
		monads: []string{"foo"},
		pairs:  map[string]string{"size": "100"},
		seq:    "ATTC",
	}

	name, size := parseAnno(seq.newSeqLinear())

	if name != seq.monads[0] {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			seq.monads[0],
			name,
		)
	}

	if sizeExpected, _ := strconv.Atoi(seq.pairs["size"]); size != sizeExpected {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			sizeExpected,
			size,
		)
	}
}
func TestParseAnnoMissingSize(t *testing.T) {
	seq := TestSequence{
		monads: []string{"foo"},
		pairs:  map[string]string{},
		seq:    "ATTC",
	}

	name, size := parseAnno(seq.newSeqLinear())

	if name != seq.monads[0] {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			seq.monads[0],
			name,
		)
	}

	sizeExpected := 1

	if size != sizeExpected {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			sizeExpected,
			size,
		)
	}
}

func TestParseAnnoUnrecognizedSize(t *testing.T) {
	seq := TestSequence{
		monads: []string{"foo"},
		pairs:  map[string]string{"size": "spam"},
		seq:    "ATTC",
	}

	name, size := parseAnno(seq.newSeqLinear())

	if name != seq.monads[0] {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			seq.monads[0],
			name,
		)
	}

	sizeExpected := 1

	if size != sizeExpected {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			sizeExpected,
			size,
		)
	}
}

func TestParseAnnoMissingName(t *testing.T) {
	seq := TestSequence{
		monads: []string{},
		pairs:  map[string]string{"size": "100"},
		seq:    "ATTC",
	}

	name, size := parseAnno(seq.newSeqLinear())

	nameExpected := "sequence"

	if name != nameExpected {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			nameExpected,
			name,
		)
	}

	if sizeExpected, _ := strconv.Atoi(seq.pairs["size"]); size != sizeExpected {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			sizeExpected,
			size,
		)
	}
}

func TestDeRep(t *testing.T) {
	seqs := []TestSequence{
		{
			monads: []string{"foo"},
			pairs:  map[string]string{"size": "100"},
			seq:    "ATTC",
		},
		{
			monads: []string{},
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

	name, size := parseAnno(<-cout)

	if nameExpected := seqs[0].monads[0]; name != nameExpected {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			nameExpected,
			name,
		)
	}

	if size != sizeExpected {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			sizeExpected,
			size,
		)
	}

	wg.Wait()
}
