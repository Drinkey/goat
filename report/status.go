package report

import (
	"io/ioutil"
	"log"
)

const (
	RUNNING = "running"
	CREATED = "created"
	DONE    = "done"
)

type Status struct {
	ID     int    `json:"-"`
	Status string `json:"status"`
	Path   Files  `json:"-"`
}

func (c *Status) set(s string) error {
	c.Status = s
	log.Printf("[taskID=%d] setting status to %s", c.ID, c.Status)
	return ioutil.WriteFile(c.Path.Status, []byte(c.Status), FILEPERM)
}

func (c Status) get() ([]byte, error) {
	log.Printf("[taskID=%d] reading status from %s", c.ID, c.Path.Status)
	return ioutil.ReadFile(c.Path.Status)
}

func (c *Status) clearCache() {
	c.Status = ""
}

func (c *Status) Get() error {
	cb, err := c.get()
	if err != nil {
		log.Printf("failed to read status from %s", c.Path.Status)
		log.Println(err)
		return err
	}
	c.Status = string(cb)
	log.Printf("[taskID=%d] got status: %s", c.ID, c.Status)
	return nil
}

func (c *Status) SetRunning() error {
	// accquire lock
	// defer unlock
	return c.set(RUNNING)
}

func (c *Status) SetCreated() error {
	// accquire lock
	// defer unlock
	return c.set(CREATED)
}

func (c *Status) SetDone() error {
	// accquire lock
	// defer unlock
	return c.set(DONE)
}
