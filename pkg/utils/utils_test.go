package utils

import (
	"os"
	"testing"
)

func TestFileExistDirReturnsFalse(t *testing.T) {
	f := "/tmp/somedir"
	os.Mkdir(f, 0777)
	defer os.Remove(f)
	r := IsFileExist(f)
	if r {
		t.Logf("expect return false but got true. %s is a DIR", f)
		t.Fail()
	}
}

func TestFileExistFileReturnsTrue(t *testing.T) {
	f := "/tmp/somefile"
	// file := os.File{}
	os.Create(f)
	defer os.Remove(f)
	r := IsFileExist(f)
	if !r {
		t.Log("expect return true but got false")
		t.Fail()
	}
}

func TestDirUpLevelUsingFullFileName(t *testing.T) {
	p := "/go/src/goat/commando/cmd.go"
	upLevel := DirUpLevel(p, -2)
	if upLevel != "/go/src/goat" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelUsingDirNameOneLevel(t *testing.T) {
	p := "/go/src/goat/commando"
	upLevel := DirUpLevel(p, -1)
	if upLevel != "/go/src/goat" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelUsingDirNameMultiLevel(t *testing.T) {
	p := "/go/src/goat/pkg/utils/utils.go"
	upLevel := DirUpLevel(p, -3)
	if upLevel != "/go/src/goat" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelPositiveLevelWillPanic(t *testing.T) {
	defer func() {
		recover()
	}()
	p := "/go/src/goat/pkg/utils/utils.go"
	_ = DirUpLevel(p, 1)
	// only fail when function is not panic
	t.Fail()
}

func TestGetCronFilePathWithoutEnviron(t *testing.T) {
	r := GetCronFilePath("jsmith")
	if r != "/var/spool/cron/jsmith" {
		t.Logf("actual: %s", r)
		t.Fail()
	}
}

func TestGetCronFilePathWithEnviron(t *testing.T) {
	os.Setenv("CRON_FILE_PATH", "/var/some/other/place")
	r := GetCronFilePath("jsmith")

	if r != "/var/some/other/place" {
		t.Logf("actual: %s", r)
		t.Fail()
	}
}

func TestLsDirReturnErrorIfPathIsNotDir(t *testing.T) {
	f := "/tmp/somefile"
	os.Create(f)
	defer os.Remove(f)
	_, err := LsDir(f)
	if err == nil {
		t.Logf("expect error but got nil")
		t.Fail()
	}
}

func TestLsDirReturnDirsUnderPathname(t *testing.T) {
	// Create two dirs and one file to run the test
	folder := "/tmp/goatcache"
	os.Mkdir(folder, 0777)
	os.Create(folder + "/testfile")
	os.Mkdir(folder+"/folder1", 0777)
	os.Mkdir(folder+"/folder2", 0777)
	defer os.RemoveAll(folder)

	dirs, err := LsDir(folder)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	f1Found, f2Found := false, false
	for _, d := range dirs {
		if d == "folder1" {
			f1Found = true
			continue
		}
		if d == "folder2" {
			f2Found = true
			continue
		}
		if d == "testfile" {
			t.Log("testfile is a file, not a dir, should not exist.")
			t.Fail()
		}
	}
	if f2Found != true || f1Found != true {
		t.Logf("dir missing: 1:%t 2:%t", f1Found, f2Found)
		t.Fail()
	}
}
