package services

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Drinkey/goat/cron"
	"github.com/Drinkey/goat/pkg/utils"
	"github.com/Drinkey/goat/report"
)

func GetTask(id int) *cron.Task {
	cron.CronTab.Parse()
	task := cron.CronTab.FindTaskByID(id)
	if task == nil {
		return nil
	}
	task.Load()
	return task
}

func IsTaskRunning(task *cron.Task) bool {
	return task.Report != nil && task.Report.Status.Status == report.RUNNING
}

func runCommand(command string, report *report.Report) ([]byte, error) {
	log.Printf("start to run command %s", command)
	report.Status.SetRunning()
	defer report.Status.SetDone()

	commandArgs := strings.Fields(command)
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to run command: %s", command)
		log.Println(err)
		return nil, err
	}
	return output, nil
}

func Execute(id int, t *cron.Task) {
	log.SetPrefix("services::Execute - ")

	log.Printf("creating report for task %d", id)
	report := report.Report{ID: id}

	report.New()

	checksum := utils.Sha256Sum(t.Schedule + " " + t.Command)
	report.Checksum.Save(checksum)

	report.Status.SetCreated()
	report.Result.SetNotRun()

	r, err := runCommand(t.Command, &report)
	if err == nil {
		report.Result.SetPass()
	} else {
		report.Result.SetFail()
	}
	log.Printf("Location: %s", report.Path.Log)
	report.Log.Save(r)
	log.Printf("-- task %d execution completed", id)
}

func GetTaskReport(id int) *report.Report {
	log.SetPrefix("services::GetTaskReport - ")
	log.Print("----")
	log.Printf("Loading task %d report", id)
	r := report.Report{ID: id}
	t := r.Load()
	if t == report.ErrTaskReportNotExist {
		return nil
	}
	return &r
}
