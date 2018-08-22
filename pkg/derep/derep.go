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
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/seqio"
)

const (
	// MinLen is the value to disable the minimal length filter
	MinLen = 0
	// MaxLen is the value to disable the maximal length filter
	MaxLen = 0
)

// Cluster is a struct that stores the name and size of an FASTA annotation.
type Cluster struct {
	linear.Seq
	Size   int
	Merged []*seq.Annotation
}

// Name returns the ID and the size of the cluster.
func (c *Cluster) Name() string {
	return fmt.Sprintf("%v;size=%d", c.ID, c.Size)
}

// PassFilter checks if the cluster can pass filter of given size.
//
// If min equals to MinLen, then the low filter is disabled. If max equals to
// MaxLen, then the high filter is disabled.
//
// Each filter is enabled when its threshold is larger than the default value.
// A cluster pass the filter when it is greater than min and smaller than max.
func (c *Cluster) PassFilter(min, max int) bool {
	passLow := (min > MinLen && c.Size >= min) || (min == MinLen)

	passHigh := (max > MaxLen && c.Size <= max) || (max == MaxLen)

	return passLow && passHigh
}

// ParseAnno parses the annotation in sequence and returns it as a cluster.
//
// The first monad is used as the name of the sequence; otherwise defaults to
// "sequence".
//
// The last key-value pair with key "size" is used as the size; otherwise
// defaults to 1.
func ParseAnno(s *linear.Seq) *Cluster {

	res := Cluster{
		Seq:    *s,
		Merged: []*seq.Annotation{&s.Annotation},
	}

	var monads []string

	pairs := make(map[string]string)

	for _, item := range strings.Split(s.ID, ";") {
		// If more than 2 then skip
		switch pair := strings.Split(item, "="); len(pair) {
		case 1:
			monads = append(monads, pair[0])
		case 2:
			pairs[pair[0]] = pair[1]
		}
	}

	res.Size = func() int {
		if val, prs := pairs["size"]; prs {
			size, err := strconv.Atoi(val)
			if err == nil && size > 0 {
				return size
			}
		}
		return 1
	}()

	res.ID = func() string {
		if len(monads) > 0 {
			return monads[0]
		}
		return "sequence"
	}()

	return &res
}

// DeRep receives a sequence from a channel and builds a map.
//
// If a sequence is in a map, DeRep parses the annotation and sums the size of
// the new sequence with the cluster in the map. if not, DeRep adds a new
// cluster into the map.
//
// After the channel in is closed, DeRep writes the map to a file.
func DeRep(in <-chan *linear.Seq, f io.Writer, min, max int, wg *sync.WaitGroup) {
	defer wg.Done()

	rep := make(map[string]*Cluster)

	for s := range in {
		c := ParseAnno(s)

		if _, prs := rep[s.String()]; !prs {
			rep[s.String()] = c
		} else {
			rep[s.String()].Size += c.Size
			rep[s.String()].Merged = append(
				rep[s.String()].Merged,
				c.Merged...,
			)
		}
	}

	w := fasta.NewWriter(f, seqio.WidthCol)

	for _, s := range rep {
		if s.PassFilter(min, max) {
			if _, err := w.Write(s); err != nil {
				log.Panicf("Error occurred during write: %s", err)
			}
		}
	}
}
