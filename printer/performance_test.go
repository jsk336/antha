// antha/printer/performance_test.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


// This file implements a simple printer performance benchmark:
// antha test -bench=BenchmarkPrint

package printer

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

var testfile *ast.File

func testprint(out io.Writer, file *ast.File) {
	if err := (&Config{TabIndent | UseSpaces, 8, 0}).Fprint(out, fset, file); err != nil {
		log.Fatalf("print error: %s", err)
	}
}

// cannot initialize in init because (printer) Fprint launches goroutines.
func initialize() {
	const filename = "testdata/parser.go"

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s", err)
	}

	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var buf bytes.Buffer
	testprint(&buf, file)
	if !bytes.Equal(buf.Bytes(), src) {
		log.Fatalf("print error: %s not idempotent", filename)
	}

	testfile = file
}

func BenchmarkPrint(b *testing.B) {
	if testfile == nil {
		initialize()
	}
	for i := 0; i < b.N; i++ {
		testprint(ioutil.Discard, testfile)
	}
}