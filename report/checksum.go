package report

import (
	"io/ioutil"
	"log"
)

type Checksum struct {
	ID        int    `json:"-"`
	Sha256sum string `json:"sha256sum"`
	Path      Files  `json:"-"`
}

func (c *Checksum) set(r string) error {
	c.Sha256sum = r
	log.Printf("[taskID=%d] setting checksum to %s", c.ID, c.Sha256sum)
	return ioutil.WriteFile(c.Path.Checksum, []byte(c.Sha256sum), FILEPERM)
}

func (c Checksum) get() ([]byte, error) {
	log.Printf("[taskID=%d] reading checksum from %s", c.ID, c.Path.Checksum)
	return ioutil.ReadFile(c.Path.Checksum)
}

func (c *Checksum) clearCache() {
	c.Sha256sum = ""
}

func (c *Checksum) Save(sum string) error {
	return c.set(sum)
}

// Get reads sha256sum from file and set to Sha256sum
func (c *Checksum) Get() error {
	cb, err := c.get()
	if err != nil {
		log.Printf("failed to read Sha256sum from %s", c.Path.Checksum)
		log.Println(err)
		return err
	}
	c.Sha256sum = string(cb)
	log.Printf("[taskID=%d] got Sha256sum: %s", c.ID, c.Sha256sum)
	return nil
}
