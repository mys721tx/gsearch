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

// Package derep provides functions to remove duplicates and sums the sequence
// abundance.
package derep

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

// Cluster is a struct that stores the name and size of an FASTA annotation.
type Cluster struct {
	Name string
	Size int
	//seqs *linear.Seq
}

// ParseAnno parses the annotation in sequence and returns it as a cluster.
//
// The first monad is used as the name of the sequence; otherwise defaults to
// "sequence".
//
// The last key-value pair with key "size" is used as the size; otherwise
// defaults to 1.
func ParseAnno(seq *linear.Seq) *Cluster {

	var monads []string

	pairs := make(map[string]string)

	for _, item := range strings.Split(seq.ID, ";") {
		// If more than 2 then skip
		switch pair := strings.Split(item, "="); len(pair) {
		case 1:
			{
				monads = append(monads, pair[0])
			}
		case 2:
			{
				pairs[pair[0]] = pair[1]
			}
		}
	}

	var size int
	var err error

	if _, prs := pairs["size"]; prs {
		size, err = strconv.Atoi(pairs["size"])
		if err != nil {
			size = 1
		}
		if size <= 0 {
			size = 1
		}
	} else {
		size = 1
	}

	var name string

	if len(monads) > 0 {
		name = monads[0]
	} else {
		name = "sequence"
	}

	res := Cluster{Name: name, Size: size}

	return &res
}

// DeRep receives a sequence from a channel and builds a map.
//
// If a sequence is in a map, DeRep parses the annotation and sums the size of
// the new sequence with the cluster in the map. if not, DeRep adds a new
// cluster into the map.
//
// After the channel out is closed, DeRep iterates through the map and send all
// sequences to the channel out.
func DeRep(in <-chan *linear.Seq, out chan<- *linear.Seq, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	rep := make(map[string]*Cluster)

	for seq := range in {
		c := ParseAnno(seq)

		if _, prs := rep[seq.String()]; !prs {
			//rep[s] = &cluster{name: name, size: size, seqs: seq}
			rep[seq.String()] = c
		} else {
			rep[seq.String()].Size += c.Size
		}
	}

	for s, v := range rep {
		out <- linear.NewSeq(
			fmt.Sprintf("%v;size=%d", v.Name, v.Size),
			[]alphabet.Letter(s),
			alphabet.DNA,
		)
	}
}
