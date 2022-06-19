package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
	"sycret/xml"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		template := c.QueryParam("URLTemplate")
		paramId := c.QueryParam("RecordID")
		id, err := strconv.Atoi(paramId)
		if err != nil {
			return c.String(http.StatusBadRequest, "неверно указан ID пациента")
		}
		file := xml.InsertInXml(template, id)
		a := make(map[string]string)
		a["URLWord"] = "/files/" + file

		return c.JSON(http.StatusOK, a)
	})
	e.GET("/files/:fileName", func(c echo.Context) error {
		url := c.Param("fileName")
		c.Response().Header().Set("Content-Disposition", "attachment; filename="+url)
		c.Response().Header().Set("Content-Type", "application/octet-stream")
		c.Response().WriteHeader(http.StatusOK)
		return c.Inline(url, url)
	})

	fmt.Printf("Файл создан")
	e.Logger.Fatal(e.Start(":8000"))
}
