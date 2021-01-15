package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Drinkey/goat/cron"
	"github.com/Drinkey/goat/pkg/app"
	"github.com/Drinkey/goat/report"
	"github.com/Drinkey/goat/services"
)

// Ping responses pong to the client. It can be used for service probing
// @Summary Response to service probing
// @Description probing
// @Produce json
// @Success 200 {string} string "ok" "PONG"
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	a.Response(http.StatusOK, app.SUCCESS, "PONG")
}

func ListCronTasks(c *gin.Context) {

	a := app.GoatResponse{Context: c}
	if err := cron.CronTab.GetReportIDs(); err != nil {
		a.Response(http.StatusInternalServerError, app.ERROR, err.Error())
	}
	cron.CronTab.Parse()
	a.Response(http.StatusOK, app.SUCCESS, cron.CronTab)
}

func RunOneTask(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.Response(http.StatusBadRequest, app.INVALID_PARAMS, "id should be int")
		return
	}
	cron.CronTab.Parse()
	task := cron.CronTab.FindTaskByID(id)
	task.Load()
	if task.Report != nil && task.Report.Status.Status == report.RUNNING {
		m := fmt.Sprintf("task %d is already running, will not run it again, query /cron/:id for details", id)
		log.Printf(m)
		a.Response(http.StatusConflict, app.ALREADY_EXIST, m)
		return
	}
	go func() {
		log.Print("task created")
		services.Execute(id, task)
	}()
	msg := fmt.Sprintf("task %d start to run, query /cron/:id for details", id)
	a.Response(http.StatusCreated, app.CREATED, msg)
}

func GetOneTask(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.Response(http.StatusBadRequest, app.INVALID_PARAMS, "id should be int")
		return
	}

	report := services.GetTask(id)
	if report == nil {
		a.Response(http.StatusNotFound, app.NOT_FOUND,
			fmt.Sprintf("report for task %d not found", id))
		return
	}
	a.Response(http.StatusOK, app.SUCCESS, &report)
}
