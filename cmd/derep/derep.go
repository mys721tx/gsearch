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
	"bufio"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/biogo/biogo/seq/linear"

	"github.com/mys721tx/gsearch/pkg/derep"
	"github.com/mys721tx/gsearch/pkg/seqio"
)

var (
	pin = flag.String(
		"in",
		"",
		"path to the sequence fasta file, default to stdin.",
	)
	pout = flag.String(
		"out",
		"",
		"path to the output fasta file, default to stdout.",
	)
	wg sync.WaitGroup
)

func main() {
	flag.Parse()

	var fin, fout *os.File
	var err error

	if *pin == "" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(*pin)
	}

	if err != nil {
		log.Panicf("failed to open %q: %v", *pin, err)
	}

	if *pout == "" {
		fout = os.Stdin
	} else {
		fout, err = os.Create(*pout)
	}

	if err != nil {
		log.Panicf("failed to open %q: %v", *pout, err)
	}

	defer func() {
		err = fout.Close()

		if err != nil {
			log.Panicf("failed to close %q: %v", *pout, err)
		}
	}()

	w := bufio.NewWriter(fout)

	defer func() {
		err = w.Flush()

		if err != nil {
			log.Panicf("failed to flush %q: %v", *pout, err)
		}
	}()

	cin := make(chan *linear.Seq)
	cout := make(chan *linear.Seq)

	wg.Add(3)
	go seqio.ScanSeq(fin, cin, &wg)
	go derep.DeRep(cin, cout, &wg)  // TODO: handling panic
	go seqio.WriteSeq(w, cout, &wg) // TODO: handling panic
	wg.Wait()
}
