package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func startServer() {
	e := echo.New()

	e.File("/", "index.html")
	e.Static("/dist", "dist")

	e.GET("/courses", coursesIndex)
	e.GET("/semesters", semestersIndex)

	fmt.Println("Starting server on http://localhost:8000")

	e.Logger.Fatal(e.Start(":8000"))
}

func coursesIndex(c echo.Context) error {
	courses := getCoursesFromDB()
	return c.JSON(http.StatusOK, courses)
}

func semestersIndex(c echo.Context) error {
	semesters := getSemestersFromDB()
	return c.JSON(http.StatusOK, semesters)
}
