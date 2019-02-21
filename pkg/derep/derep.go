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
	"io"
	"log"

	"github.com/biogo/biogo/io/seqio/fasta"
	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/cluster"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

// DeRep receives a sequence from a channel and builds a map.
//
// If a sequence is in a map, DeRep parses the annotation and sums the size of
// the new sequence with the cluster in the map. if not, DeRep adds a new
// cluster into the map.
//
// After the channel in is closed, DeRep writes the map to a file.
func DeRep(in []*linear.Seq, f io.Writer, min, max int) {
	rep := make(map[string]*cluster.Cluster)

	for _, s := range in {
		c := cluster.ParseAnno(s)

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
