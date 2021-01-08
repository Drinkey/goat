package services

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Drinkey/goat/commando"
)

func runCommand(command string, setStatus *commando.Status) ([]byte, error) {
	log.Printf("start to run command %s", command)
	setStatus.Running()
	defer setStatus.Done()

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

func Execute(id int, cmd string) {
	log.SetPrefix("services::Execute - ")
	p := commando.Files{ID: id}
	p.Load()
	p.CreateDir()
	setStatus, setResult := commando.Status{Path: p, ID: id}, commando.Result{Path: p, ID: id}

	setStatus.Created()
	setResult.NotRun()

	r, err := runCommand(cmd, &setStatus)
	if err == nil {
		setResult.Pass()
	} else {
		setResult.Fail()
	}
	setResult.SaveLastLog(r)
}

func GetTaskStatusAndResult(id int) *commando.Report {
	log.SetPrefix("services::GetTaskStatusAndResult - ")
	report := commando.Report{ID: id}
	t := report.Load()
	if t == commando.ErrTaskReportNotExist {
		return nil
	}
	return &report
}
