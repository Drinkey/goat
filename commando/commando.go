package commando

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Drinkey/goat/pkg/utils"
)

const (
	RUNNING  = "running"
	CREATED  = "created"
	DONE     = "done"
	PASS     = "pass"
	FAIL     = "fail"
	NOTRUN   = "not run"
	FILEPERM = 0755
)

var ErrTaskReportNotExist = errors.New("Task Report does not exist")

// var CACHEDIR = utils.GetCacheDir()

type Report struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Result  string `json:"result"`
	LogPath string `json:"log_path"`
}

func (r *Report) Load() error {
	p := Files{ID: r.ID}
	p.Load()
	if !utils.IsDirExist(p.Dir) {
		log.Printf("dir not exist, skip : %s", p.Dir)
		return ErrTaskReportNotExist
	}

	status, result := Status{ID: r.ID, Path: p}, Result{ID: r.ID, Path: p}

	if err := status.Get(); err != nil {
		log.Printf("unable to get status")
		return err
	}
	if err := result.Get(); err != nil {
		log.Printf("unable to get result")
		return err
	}
	r.Status = status.Status
	r.Result = result.Result
	r.LogPath = p.LastLog
	return nil
}

// Files is an aggregation of commando files
type Files struct {
	ID                          int
	Dir                         string
	Status, LastResult, LastLog string
}

// Load initiates cache dir and set properties values
func (c *Files) Load() {
	c.Dir = fmt.Sprintf("%s/%d", utils.GetCacheDir(), c.ID)

	c.Status = c.Dir + "/status"
	c.LastResult = c.Dir + "/lastresult"
	c.LastLog = c.Dir + "/lastlog"
}

func (c Files) CreateDir() {
	// create cache dir if not exists
	if !utils.IsDirExist(c.Dir) {
		log.Printf("report for task %d is not exist, creating the folder %s", c.ID, c.Dir)
		os.MkdirAll(c.Dir, FILEPERM)
	}
}

type Result struct {
	ID     int    `json:"-"`
	Result string `json:"last_result"`
	Path   Files  `json:"-"`
}

func (c *Result) set(r string) error {
	c.Result = r
	log.Printf("[taskID=%d] setting lastresult to %s", c.ID, c.Result)
	return ioutil.WriteFile(c.Path.LastResult, []byte(c.Result), FILEPERM)
}

func (c *Result) Pass() error {
	return c.set(PASS)
}

func (c *Result) Fail() error {
	return c.set(FAIL)
}

func (c *Result) NotRun() error {
	// if file not exist, not to not run, otherwise do not update last result
	if utils.IsFileExist(c.Path.LastResult) {
		log.Print("LastResult already exist, not overwrite it")
		return nil
	}
	return c.set(NOTRUN)
}

func (c Result) SaveLastLog(output []byte) error {
	log.Printf("[taskID=%d] save lastlog to %s", c.ID, c.Path.LastLog)
	return ioutil.WriteFile(c.Path.LastLog, output, FILEPERM)
}

func (c *Result) Get() error {
	log.Printf("[taskID=%d] reading last result from %s", c.ID, c.Path.LastResult)
	cb, err := ioutil.ReadFile(c.Path.LastResult)
	if err != nil {
		log.Printf("failed to read last result from %s", c.Path.LastResult)
		return err
	}
	c.Result = string(cb)
	log.Printf("[taskID=%d] got last result: %s", c.ID, c.Result)
	return nil
}

type Status struct {
	ID     int    `json:"-"`
	Status string `json:"current_status"`
	Path   Files  `json:"-"`
}

func (c *Status) set(s string) error {
	c.Status = s
	log.Printf("[taskID=%d] setting status to %s", c.ID, c.Status)
	return ioutil.WriteFile(c.Path.Status, []byte(c.Status), FILEPERM)
}

func (c *Status) Get() error {
	log.Printf("[taskID=%d] reading status from %s", c.ID, c.Path.Status)
	cb, err := ioutil.ReadFile(c.Path.Status)
	if err != nil {
		log.Printf("failed to read status from %s", c.Path.Status)
		log.Println(err)
		return err
	}
	c.Status = string(cb)
	log.Printf("[taskID=%d] got status: %s", c.ID, c.Status)
	return nil
}

func (c *Status) Running() error {
	// accquire lock
	// defer unlock
	return c.set(RUNNING)
}

func (c *Status) Created() error {
	// accquire lock
	// defer unlock
	return c.set(CREATED)
}

func (c *Status) Done() error {
	// accquire lock
	// defer unlock
	return c.set(DONE)
}

func init() {
	log.SetPrefix("commando - ")
}
