package report

import (
	"os"
	"testing"
)

// var testid = 123

func TestStatusTransition(t *testing.T) {
	os.Setenv("GOAT_CACHE_DIR", "/tmp/goattest")
	var p = Files{ID: testid}
	p.Init()
	p.CreateDir()
	stat := Status{ID: testid, Path: p}
	// set status to Created, then clear the status cache in current struct
	stat.SetCreated()
	stat.clearCache()

	// then read the status again
	if err := stat.Get(); err != nil {
		t.Fail()
	}
	if stat.Status != CREATED {
		t.Fail()
	}

	stat.SetDone()
	stat.clearCache()

	// then read the status again
	if err := stat.Get(); err != nil {
		t.Fail()
	}
	if stat.Status != DONE {
		t.Fail()
	}
}
