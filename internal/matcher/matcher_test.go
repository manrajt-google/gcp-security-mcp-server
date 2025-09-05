package matcher

import (
	"testing"

	"github.com/google/osv-scalibr/extractor"
	"github.com/google/osv-scalibr/purl"
)

func TestNewLocalMatcher(t *testing.T) {
	matcher := NewLocalMatcher()
	if matcher == nil {
		t.Errorf("NewLocalMatcher() returned nil")
	}
}

func TestMatchVulnerabilities(t *testing.T) {
	matcher := NewLocalMatcher()

	if matcher == nil {
		t.Errorf("NewLocalMatcher() returned nil")
	}

	p := &extractor.Package{
		Name:     "github.com/google/osv-scalibr",
		Version:  "v0.3.2",
		PURLType: purl.TypeGolang,
	}

	invs := []*extractor.Package{p}

	matcher.MatchVulnerabilities(t.Context(), invs)
}
