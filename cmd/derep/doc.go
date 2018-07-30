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

/*
DeRep removes duplications from FASTA sequences and sums the sequence abundance.

DeRep scans through a FASTA file and builds a map using the sequence as key.

DeRep requires the FASTA header of a sequence to be a semicolon delimited list:
	> NM_000518;HBB;size=628;organism=9606
	ACATTTGCTT
	...

DeRep uses the first monad, a field that does not have an equal sign, as the
name of the sequence. If such monad does not exist in the header, the name of
the sequence defaults to "sequence".

DeRep uses the last key-value pair, a field that has an equal sign, with "size"
as key and an integer as value for the abundance of that sequence. If such pair
does not exist in the header, the size of the sequence defaults to 1.

Usage:
	derep [flags]

The flags are:
	-in string
		path to the sequence FASTA file, default to stdin.
	-max int
		maximal abundance of a sequence, default to 0.
	-min int
		minimal abundance of a sequence, default to 0.
	-out string
		path to the output FASTA file, default to stdout.

Example:
	derep -in short.fasta -out merged.fasta
	gunzip -c compress_seq.fasta.gz | derep -out merged.fasta
	derep -in short.fasta | grep ">" | sort
*/
package main

// BUG(mys721tx): Currently DeRep does not preserve original FASTA header other
// than the name and abundance of the sequence.
