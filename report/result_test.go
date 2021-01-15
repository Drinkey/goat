package report

import (
	"os"
	"testing"
)

var testid = 123

func TestResultTransition(t *testing.T) {
	os.Setenv("GOAT_CACHE_DIR", "/tmp/goattest")
	var p = Files{ID: testid}
	p.Init()
	p.CreateDir()

	result := Result{ID: testid, Path: p}

	if err := result.SetPass(); err != nil {
		t.Fail()
	}
	result.clearCache()

	if err := result.Get(); err != nil {
		t.Fail()
	}
	if result.Result != pass {
		t.Fail()
	}

	if err := result.SetNotRun(); err != nil {
		t.Fail()
	}
	result.clearCache()

	if err := result.Get(); err != nil {
		t.Fail()
	}
	if result.Result != notrun {
		t.Fail()
	}
}
