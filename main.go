package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	s := echo.New()
	s.HideBanner = true
	s.Use(middleware.Logger())

	// TODO: squeue -> tasks
	s.GET("/api/v1/slurm/squeue", squeue)
	s.POST("/api/v1/slurm/sbatch", sbatch)
	s.GET("/api/v1/slurm/scontrol", scontrol)
	s.DELETE("/api/v1/slurm/scancel", scancel)

	s.Logger.Error(s.Start(":1888"))
}

func squeue(c echo.Context) error {
	r, err := runSqueue()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if r == nil {
		return c.String(http.StatusOK, "[]")
	}
	return c.JSON(http.StatusOK, r)
}

func sbatch(c echo.Context) error {
	workDir := c.QueryParam("work_dir")
	arg := c.QueryParam("arg")
	file := c.QueryParam("file")
	r, err := runSbatch(workDir, arg, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, r)
}

func scancel(c echo.Context) error {
	id := c.QueryParam("job_id")
	err := runScancel(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Canceled")
}

func scontrol(c echo.Context) error {
	id := c.QueryParam("job_id")
	r, err := runScontrol(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, r)
}
