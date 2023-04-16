package main

//go:generate go build -o=../../bin/staticlint

import (
	"strings"

	"github.com/ksusonic/go-devops-mon/pkg/exitcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

var analyzers = []*analysis.Analyzer{
	// report mismatches between assembly files and Go declarations
	asmdecl.Analyzer,
	// check for useless assignments
	assign.Analyzer,
	// check for common mistakes using the sync/atomic package
	atomic.Analyzer,
	// checks for non-64-bit-aligned arguments to sync/atomic functions
	atomicalign.Analyzer,
	// check for common mistakes involving boolean operators
	bools.Analyzer,
	// check that +build tags are well-formed and correctly located
	buildtag.Analyzer,
	// detect some violations of the cgo pointer passing rules. Not working in Arcadia for now
	// cgocall.Analyzer,
	// check for unkeyed composite literals
	composite.Analyzer,
	// check for locks erroneously passed by value
	copylock.Analyzer,
	// check for the use of reflect.DeepEqual with error values
	deepequalerrors.Analyzer,
	// check that the second argument to errors.As is a pointer to a type implementing error
	errorsas.Analyzer,
	// check for mistakes using HTTP responses
	httpresponse.Analyzer,
	// check references to loop variables from within nested functions
	loopclosure.Analyzer,
	// check cancel func returned by context.WithCancel is called
	lostcancel.Analyzer,
	// check for useless comparisons between functions and nil
	nilfunc.Analyzer,
	// inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons
	nilness.Analyzer,
	// check consistency of Printf format strings and arguments
	printf.Analyzer,
	// check for possible unintended shadowing of variables EXPERIMENTAL
	// shadow.Analyzer,
	// check for shifts that equal or exceed the width of the integer
	shift.Analyzer,
	// check signature of methods of well-known interfaces
	stdmethods.Analyzer,
	// check that struct field tags conform to reflect.StructTag.Get
	structtag.Analyzer,
	// check for common mistaken usages of tests and examples
	tests.Analyzer,
	// report passing non-pointer or non-interface values to unmarshal
	unmarshal.Analyzer,
	// check for unreachable code
	unreachable.Analyzer,
	// check for invalid conversions of uintptr to unsafe.Pointer
	unsafeptr.Analyzer,
	// check for unused results of calls to some functions
	unusedresult.Analyzer,
	// check for unused writes
	unusedwrite.Analyzer,
	// check for string(int) conversions
	stringintconv.Analyzer,
	// check for impossible interface-to-interface type assertions
	ifaceassert.Analyzer,

	// Checks os.Exit from main
	exitcheck.Analyzer,
}

var additionalStaticAnalyzers = map[string]bool{
	"ST1023": true, // Redundant type in variable declaration
	"QF1003": true, // Convert if/else-if chain to tagged switch
}

func main() {
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		} else if _, ok := additionalStaticAnalyzers[v.Analyzer.Name]; ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
