package report

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Drinkey/goat/pkg/utils"
)

// ErrTaskReportNotExist means the report of specific task id does not exist
var ErrTaskReportNotExist = errors.New("Task Report does not exist")

const (
	FILEPERM = 0755
)

// Files is an aggregation of commando files
type Files struct {
	ID                  int
	Dir, Checksum       string
	Status, Result, Log string
}

// Init initiates cache dir and set properties values
func (c *Files) Init() {
	c.Dir = fmt.Sprintf("%s/%d", utils.GetCacheDir(), c.ID)

	c.Status = c.Dir + "/status"
	c.Result = c.Dir + "/result"
	c.Log = c.Dir + "/log"
	c.Checksum = c.Dir + "/SHASUM"
}

// CreateDir creates cache dir if not exists
func (c Files) CreateDir() {
	if !utils.IsDirExist(c.Dir) {
		log.Printf("report for task %d is not exist, creating the folder: %s", c.ID, c.Dir)
		os.MkdirAll(c.Dir, FILEPERM)
	}
}

func init() {
	log.SetPrefix("Report - ")
}

// Report is composited of Result and Status
type Report struct {
	ID       int      `json:"id"`
	Status   Status   `json:"status"`
	Result   Result   `json:"result"`
	Log      Log      `json:"log"`
	Checksum Checksum `json:"checksum"`
	Path     Files    `json:"-"`
}

// New creates new report without checking if dir existence
func (r *Report) New() error {
	return r.initialize(true, false)
}

func (r *Report) initialize(createDir bool, errOnDirNotExit bool) error {
	p := Files{ID: r.ID}
	p.Init()

	if createDir {
		p.CreateDir()
	}
	if errOnDirNotExit && !utils.IsDirExist(p.Dir) {
		log.Printf("dir not exist, skip : %s", p.Dir)
		return ErrTaskReportNotExist
	}
	r.Path = p
	r.Status = Status{ID: r.ID, Path: r.Path}
	r.Result = Result{ID: r.ID, Path: r.Path}
	r.Checksum = Checksum{ID: r.ID, Path: r.Path}
	r.Log = Log{ID: r.ID, Path: r.Path}
	return nil
}

func (r *Report) readStatus() error {
	if err := r.Status.Get(); err != nil {
		log.Printf("unable to get status")
		return err
	}
	return nil
}

func (r *Report) readResult() error {

	if err := r.Result.Get(); err != nil {
		log.Printf("unable to get result")
		return err
	}
	return nil
}

func (r *Report) readChecksum() error {

	if err := r.Checksum.Get(); err != nil {
		log.Printf("unable to get checksum")
		return err
	}
	return nil
}

func (r *Report) readLog() error {
	if err := r.Log.Get(); err != nil {
		log.Printf("unable to get log")
		return err
	}
	return nil
}

// Load reads exist report files and get values, returns error if report not exist
func (r *Report) Load() error {
	if err := r.initialize(false, true); err != nil {
		return err
	}
	if err := r.readStatus(); err != nil {
		return err
	}
	if err := r.readResult(); err != nil {
		return err
	}
	if err := r.readChecksum(); err != nil {
		return err
	}

	if err := r.readLog(); err != nil {
		return err
	}

	return nil
}
