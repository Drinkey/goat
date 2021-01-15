package report

import (
	"io/ioutil"
	"log"
)

type Log struct {
	ID       int    `json:"-"`
	Content  string `json:"content"`
	Path     Files  `json:"-"`
	Location string `json:"log_path"`
}

func (c *Log) Save(output []byte) error {
	c.Location = c.Path.Log
	log.Printf("[taskID=%d] save log to %s", c.ID, c.Location)
	return ioutil.WriteFile(c.Path.Log, output, FILEPERM)
}

func (c *Log) Get() error {
	c.Location = c.Path.Log
	log.Printf("[taskID=%d] reading log from %s", c.ID, c.Location)
	cb, err := ioutil.ReadFile(c.Location)
	if err != nil {
		log.Println(err)
		return err
	}

	c.Content = string(cb)
	log.Println(c.Content)
	return nil
}
