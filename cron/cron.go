package cron

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Drinkey/goat/pkg/utils"
	"github.com/Drinkey/goat/report"
	"github.com/gorhill/cronexpr"
)

// CronTab is a global singleton object
var CronTab Cron

// Task represents one task of crontab
type Task struct {
	ID        int            `json:"id"`
	Title     string         `json:"title"`
	Schedule  string         `json:"schedule"`
	Command   string         `json:"command"`
	NextRun   time.Time      `json:"next_run"`
	Report    *report.Report `json:"report,omitempty"`
	Checksum  string         `json:"checksum"`
	IsChanged bool           `json:"is_changed_since_lastrun"`
}

// Load reads the execution report of current task
func (t *Task) Load() {
	r := &report.Report{ID: t.ID}
	log.Printf("reading task %d report", r.ID)
	err := r.Load()
	if err != nil {
		log.Printf("something wrong when read task report")
		log.Println(err)
		return
	}
	t.Report = r
	log.Printf("Current checksum %s, previous report checksum %s", t.Checksum, r.Checksum.Sha256sum)
	t.IsChanged = !(t.Checksum == r.Checksum.Sha256sum)
	log.Printf("IsChanged=%t", t.IsChanged)
}

// Cron represents the info of crontab
type Cron struct {
	Host      string       `json:"host"`
	User      string       `json:"user"`
	TimeZone  string       `json:"timezone"`
	Count     int          `json:"task_count"`
	Tasks     []*Task      `json:"tasks"`
	reportIDs map[int]bool `json:"-"`
	File      string       `json:"-"`
}

func (c Cron) parseLine(index int, line string, lastline string) *Task {
	e := strings.Fields(line)
	sched := strings.Join(e[:5], " ")
	title := ""
	// Title must be the lastline with format "# xxx", xxx will be title
	if strings.Index(lastline, "# ") == 0 {
		title = strings.Replace(lastline, "# ", "", 1)
	} else {
		log.Printf("unrecognized title format: %s", lastline)
	}
	return &Task{
		ID:       index,
		Title:    title,
		Schedule: sched,
		NextRun:  cronexpr.MustParse(sched).Next(time.Now()),
		Command:  strings.Join(e[5:], " "),
		Checksum: utils.Sha256Sum(line),
	}
}

func (c Cron) parseCronTab(content []byte) []*Task {
	t := []*Task{}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Index(line, "#") == 0 || len(line) == 0 {
			// skip useless lines
			continue
		}
		index := len(t) + 1
		var task *Task
		if i > 0 {
			task = c.parseLine(index, line, lines[i-1])
		} else {
			task = c.parseLine(index, line, "")
		}

		if _, ok := c.reportIDs[len(t)+1]; ok {
			log.Printf("report for task %d exist, read it", index)
			task.Load()
		} else {
			task.Report = nil
		}
		t = append(t, task)
	}
	return t
}

// Parse parses cron file into Cron struct
func (c *Cron) Parse() {
	log.Printf("Parsing cron file %s", c.File)
	content, err := ioutil.ReadFile(c.File)
	if err != nil {
		log.Printf("failed to read cron file %s: %s", c.File, err.Error())
		log.Print("failover to use command to get crontab")
		content, err = exec.Command("crontab", "-l").Output()
		if err != nil {
			log.Printf("crontab -l returns error: %s", err.Error())
			log.Panic("neither command or reading file can retrieve cron jobs")
		}
	}

	c.Tasks = c.parseCronTab(content)
	c.Count = len(c.Tasks)
}

// GetReportIDs walks through cache dir and parse existing execution report, save it to reportIDs
func (c *Cron) GetReportIDs() error {
	cacheDir := utils.GetCacheDir()
	log.Printf("GOAT_CACHE_DIR=%s", cacheDir)
	dirs, err := utils.LsDir(cacheDir)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Got tasks as following:")
	log.Println(dirs)

	// init reportIDs, it'll only initialize once
	c.reportIDs = map[int]bool{}
	for _, t := range dirs {
		reportID, err := strconv.Atoi(t)

		if err != nil {
			log.Printf("failed to covert report id %s to int", t)
			log.Println(err)
			continue
		}
		if reportID <= 0 {
			log.Printf("got illegal report id %d", reportID)
			continue
		}
		c.reportIDs[reportID] = true
	}
	return nil
}

// SetTimeZone sets timezone info to c.TimeZone. If param tz exists, set c.TimeZone to tz
func (c *Cron) SetTimeZone(tz ...string) {
	if len(tz) != 0 {
		c.TimeZone = tz[0]
		return
	}
	t := time.Now()
	zone, _ := t.Local().Zone()

	c.TimeZone = zone
	log.Printf("Set Timezone to %s", c.TimeZone)
}

// SetHost sets host info to c.Host. If param hostname exists, set c.Host to hostname
func (c *Cron) SetHost(hostname ...string) {
	if len(hostname) != 0 {
		c.Host = hostname[0]
		return
	}
	h, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname")
	}

	c.Host = h
	log.Printf("Set Host to %s", c.Host)
}

// FindTaskByID returns the task pointer with specified task id
func (c Cron) FindTaskByID(id int) *Task {
	for _, task := range c.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func init() {
	log.SetPrefix("cron::init - ")
	u := utils.GetWhoAmI()
	cronFile := utils.GetCronFilePath(u)

	CronTab = Cron{User: u, File: cronFile}
	CronTab.SetHost()
	CronTab.SetTimeZone()

	// create cache dir if not exist
	cache := utils.GetCacheDir()
	if !utils.IsDirExist(cache) {
		log.Printf("cache dir %s does not exist, create it", cache)
		if err := os.MkdirAll(cache, 0700); err != nil {
			log.Panicf("failed to create cache dir %s", cache)
		}
	}
}
