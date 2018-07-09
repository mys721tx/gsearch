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
	"testing"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

type TestSequence struct {
	name string
	size int
	seq  string
}

func (t *TestSequence) newSeqLinear() *linear.Seq {
	return linear.NewSeq(
		fmt.Sprintf("%s;size=%d", t.name, t.size),
		[]alphabet.Letter(t.seq),
		alphabet.DNA,
	)
}

func TestParseAnno(t *testing.T) {
	seq := TestSequence{name: "foo", size: 100, seq: "ATTC"}

	name, size := parseAnno(seq.newSeqLinear())

	if name != seq.name {
		t.Errorf(
			"expecting name %s, parseAnno returns %s",
			seq.name,
			name,
		)
	}

	if size != seq.size {
		t.Errorf(
			"expecting size %d, parseAnno returns %d",
			seq.size,
			size,
		)
	}
}
