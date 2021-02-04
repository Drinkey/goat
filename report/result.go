package report

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Drinkey/goat/pkg/utils"
)

const (
	pass   = "pass"
	fail   = "fail"
	notrun = "not run"
)

// Result is the object represent the task execution result
type Result struct {
	ID     int       `json:"-"`
	Result string    `json:"result"`
	Time   time.Time `json:"time"`
	Path   Files     `json:"-"`
}

func (c *Result) set(r string) error {
	c.Result = r
	log.Printf("[taskID=%d] setting result to %s", c.ID, c.Result)
	return ioutil.WriteFile(c.Path.Result, []byte(c.Result), FILEPERM)
}

func (c Result) get() ([]byte, error) {
	log.Printf("[taskID=%d] reading result from %s", c.ID, c.Path.Result)
	return ioutil.ReadFile(c.Path.Result)
}

func (c *Result) clearCache() {
	c.Result = ""
}

// SetPass sets result to pass
func (c *Result) SetPass() error {
	return c.set(pass)
}

// SetFail sets result to fail
func (c *Result) SetFail() error {
	return c.set(fail)
}

// SetNotRun sets result to not run
func (c *Result) SetNotRun() error {
	// if file not exist, not to not run, otherwise do not update last result
	if utils.IsFileExist(c.Path.Result) {
		log.Print("Result already exist, not overwrite it")
		return nil
	}
	return c.set(notrun)
}

func (c *Result) getModificationTime() error {
	f, err := os.Open(c.Path.Result)
	if err != nil {
		log.Printf("failed to open last result for stat")
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		log.Printf("failed to stat last result for last update time")
		return err
	}
	t := stat.ModTime().Format(time.RFC3339)
	c.Time, err = time.Parse(time.RFC3339, t)
	if err != nil {
		log.Printf("[taskID=%d] parse modification time=%s error: %s", c.ID, c.Time, err.Error())
		return err
	}
	log.Printf("[taskID=%d] got last result modification time: %s", c.ID, c.Time)
	return nil
}

func (c *Result) getResult() error {
	log.Printf("[taskID=%d] reading last result from %s", c.ID, c.Path.Result)
	cb, err := ioutil.ReadFile(c.Path.Result)
	if err != nil {
		log.Printf("failed to read last result from %s", c.Path.Result)
		return err
	}
	c.Result = string(cb)
	log.Printf("[taskID=%d] got last result: %s", c.ID, c.Result)
	return nil
}

// Get sets last result and last modification time to Result
func (c *Result) Get() error {
	if err := c.getResult(); err != nil {
		return err
	}

	if err := c.getModificationTime(); err != nil {
		return err
	}

	return nil
}
