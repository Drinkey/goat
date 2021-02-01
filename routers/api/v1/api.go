package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Drinkey/goat/cron"
	"github.com/Drinkey/goat/pkg/app"
	"github.com/Drinkey/goat/services"
)

// Ping responses pong to the client. It can be used for service probing
// @Summary Response to service probing
// @Description probing
// @Produce json
// @Success 200 {string} string ok PONG
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	a.Response(http.StatusOK, app.SUCCESS, "PONG")
}

// ListCronTasks List all cron tasks of the host with execution status
// @Summary List all cron tasks
// @Description  List all cron tasks of the running host with execution status
// @Tags Cron
// @Produce json
// @Success 200 {object} cron.Cron cronjobs
// @Failure 500 {string} string "error message"
// @Router /api/v1/cron [get]
func ListCronTasks(c *gin.Context) {

	a := app.GoatResponse{Context: c}
	if err := cron.CronTab.GetReportIDs(); err != nil {
		a.Response(http.StatusInternalServerError, app.ERROR, err.Error())
		return
	}
	cron.CronTab.Parse()
	a.Response(http.StatusOK, app.SUCCESS, cron.CronTab)
}

// RunOneTask runs a task by specified ID
// @Summary Run a task by specified ID
// @Description Run a task by specified ID for once
// @Tags Cron
// @Produce json
// @Param  id path int true "Task ID"
// @Success 200 {object} string "Start task success"
// @Failure 400 {string} string "Invalid Request"
// @Failure 409 {string} string "Task already running"
// @Router /api/v1/cron/{id} [post]
func RunOneTask(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.Response(http.StatusBadRequest, app.INVALID_PARAMS, "id should be int")
		return
	}
	task := services.GetTask(id)
	if task == nil {
		m := fmt.Sprintf("failed to find task by task id %d, task not exist", id)
		log.Print(m)
		a.Response(http.StatusNotFound, app.NOT_FOUND, m)
		return
	}

	if services.IsTaskRunning(task) {
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

// GetOneTask responses task execution status of task by specified ID
// @Summary Get execution status
// @Description Get execution status of task by specified ID
// @Tags Cron
// @Produce json
// @Param  id path int true "Task ID"
// @Success 200 {object} report.Report "Get the task success"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task report not found"
// @Router /api/v1/cron/{id} [get]
func GetOneTask(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.Response(http.StatusBadRequest, app.INVALID_PARAMS, "id should be int")
		return
	}

	report := services.GetTaskReport(id)
	if report == nil {
		a.Response(http.StatusNotFound, app.NOT_FOUND,
			fmt.Sprintf("report for task %d not found", id))
		return
	}
	a.Response(http.StatusOK, app.SUCCESS, &report)
}
