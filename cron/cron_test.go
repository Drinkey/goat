package cron

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Drinkey/goat/pkg/utils"
)

func getTestFilePath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(pwd)
	rootPath := utils.DirUpLevel(pwd, -1)
	return fmt.Sprintf("%s/tests", rootPath)
}
func TestCronParseCrontabFileToFields(t *testing.T) {
	testFileDir := getTestFilePath()
	u := "automation"
	cronFile := fmt.Sprintf("%s/%s", testFileDir, u)

	cron := Cron{User: u, File: cronFile}
	cron.SetHost("Test")
	cron.SetTimeZone("UTC")
	cron.Parse()

	if cron.Host != "Test" {
		t.Fail()
	}
	if cron.User != u {
		t.Fail()
	}
	if len(cron.Tasks) == 0 {
		t.Fail()
	}
	firstTask := cron.Tasks[0]
	if firstTask.ID != 1 || firstTask.Schedule != "30 18 * * *" {
		t.Fail()
	}
}

func TestCronParseUserNotExist(t *testing.T) {
	defer func() {
		recover()
	}()
	testFileDir := getTestFilePath()
	u := "automationcc"
	cronFile := fmt.Sprintf("%s/%s", testFileDir, u)

	cron := Cron{User: u, File: cronFile}
	cron.Parse()
	t.Fail()
}
