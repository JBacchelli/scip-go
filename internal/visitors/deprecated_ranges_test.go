package visitors

import (
	"slices"
	"testing"

	"github.com/scip-code/scip/bindings/go/scip"
)

func TestDeprecatedRangeSlice(t *testing.T) {
	singleLine := scip.Range{
		Start: scip.Position{Line: 5, Character: 2},
		End:   scip.Position{Line: 5, Character: 9},
	}
	if got, want := deprecatedRangeSlice(singleLine), []int32{5, 2, 9}; !slices.Equal(got, want) {
		t.Fatalf("single-line: got %v, want %v", got, want)
	}

	multiLine := scip.Range{
		Start: scip.Position{Line: 5, Character: 2},
		End:   scip.Position{Line: 7, Character: 1},
	}
	if got, want := deprecatedRangeSlice(multiLine), []int32{5, 2, 7, 1}; !slices.Equal(got, want) {
		t.Fatalf("multi-line: got %v, want %v", got, want)
	}
}

func TestBackfillDeprecatedRanges(t *testing.T) {
	r := scip.Range{Start: scip.Position{Line: 1, Character: 0}, End: scip.Position{Line: 1, Character: 4}}
	occ := &scip.Occurrence{TypedRange: r.AsTypedRange(), Symbol: "x"}
	withEnclosing := &scip.Occurrence{
		TypedRange:          r.AsTypedRange(),
		TypedEnclosingRange: r.AsTypedEnclosingRange(),
		Symbol:              "y",
	}

	// nil entries must be tolerated.
	backfillDeprecatedRanges([]*scip.Occurrence{occ, nil, withEnclosing})

	if got := occ.GetRange(); !slices.Equal(got, []int32{1, 0, 4}) {
		t.Fatalf("deprecated range not backfilled: %v", got)
	}
	// The typed range is a dual-write, not a replacement -- it must survive.
	if _, ok := occ.SourceRange(); !ok {
		t.Fatal("typed range lost after backfill")
	}
	if got := withEnclosing.GetEnclosingRange(); !slices.Equal(got, []int32{1, 0, 4}) {
		t.Fatalf("deprecated enclosing range not backfilled: %v", got)
	}
}

func TestBackfillDeprecatedRangesPreservesExisting(t *testing.T) {
	// An occurrence that already carries a deprecated range must not be rewritten.
	occ := &scip.Occurrence{Range: []int32{9, 9, 9}, Symbol: "z"}
	backfillDeprecatedRanges([]*scip.Occurrence{occ})
	if got := occ.GetRange(); !slices.Equal(got, []int32{9, 9, 9}) {
		t.Fatalf("existing deprecated range overwritten: %v", got)
	}
}
