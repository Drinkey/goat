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
// @Success 200 {string} string "ok" "PONG"
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	a.Response(http.StatusOK, app.SUCCESS, "PONG")
}

func ListCronTasks(c *gin.Context) {

	a := app.GoatResponse{Context: c}
	cron.CronTab.Parse()
	cron.CronTab.ParseReport()
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
	go func() {
		log.Print("task created")
		services.Execute(id, task.Command)
	}()
	a.Response(http.StatusCreated, app.CREATED, task)
}

func GetOneTask(c *gin.Context) {
	a := app.GoatResponse{Context: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		a.Response(http.StatusBadRequest, app.INVALID_PARAMS, "id should be int")
		return
	}

	report := services.GetTaskStatusAndResult(id)
	if report == nil {
		a.Response(http.StatusNotFound, app.NOT_FOUND,
			fmt.Sprintf("report for task %d not found", id))
		return
	}
	a.Response(http.StatusOK, app.SUCCESS, &report)
}
