package cron

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Drinkey/goat/commando"
	"github.com/Drinkey/goat/pkg/utils"
	"github.com/gorhill/cronexpr"
)

const cronFileDir = "/var/spool/cron/crontabs"

// CronTab is a global singleton object
var CronTab Cron

// Task represents one task of crontab
type Task struct {
	ID               int       `json:"id"`
	Schedule         string    `json:"schedule"`
	ScheduleReadable string    `json:"schedule_readable"`
	Command          string    `json:"command"`
	LastResult       string    `json:"last_result"`
	NextRun          time.Time `json:"next_run"`
	Status           string    `json:"status"`
}

// Load reads the execution report of current task
func (t *Task) Load() {
	report := commando.Report{ID: t.ID}
	log.Printf("reading task %d report", report.ID)
	report.Load()
	t.LastResult = report.Result
	t.Status = report.Status
}

// Cron represents the info of crontab
type Cron struct {
	Host     string  `json:"host"`
	TimeZone string  `json:"timezone"`
	Count    int     `json:"task_count"`
	Tasks    []*Task `json:"tasks"`
	User     string  `json:"user"`
	File     string  `json:"-"`
}

func (c Cron) parseLine(index int, line string) *Task {
	e := strings.Fields(line)
	sched := strings.Join(e[:5], " ")
	return &Task{
		ID:       index,
		Schedule: sched,
		NextRun:  cronexpr.MustParse(sched).Next(time.Now()),
		Command:  strings.Join(e[5:], " "),
	}
}

func (c Cron) parseCronTab(content []byte) []*Task {
	t := []*Task{}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Index(line, "#") == 0 || len(line) == 0 {
			// skip useless lines
			continue
		}
		task := c.parseLine(len(t)+1, line)
		t = append(t, task)
	}
	return t
}

// Parse parses cron file into Cron struct
func (c *Cron) Parse() {
	log.Printf("Parsing file %s", c.File)
	content, err := ioutil.ReadFile(c.File)
	if err != nil {
		log.Panicf("failed to read file %s: %s", c.File, err.Error())
	}
	c.Tasks = c.parseCronTab(content)
	c.Count = len(c.Tasks)
}

// ParseReport walks through cache dir and parse existing execution report
func (c *Cron) ParseReport() {
	cacheDir := utils.GetCacheDir()
	log.Printf("GOAT_CACHE_DIR=%s", cacheDir)
	dirs, err := utils.LsDir(cacheDir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Got tasks as following:")
	log.Println(dirs)
	for _, t := range dirs {
		taskID, err := strconv.Atoi(t)
		if err != nil {
			log.Printf("failed to covert %s to int", t)
			log.Println(err)
			continue
		}
		c.Tasks[taskID-1].Load()
	}
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
	// CronTab.Parse()
	// CronTab.ParseReport()
}
