package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	// If there are 2 or more command line
	// arguments then check what the second
	// one is.
	// i.e. go run *.go <os.Args[1]>
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "import":
			importCourses()
		default:
			runServer()
		}
	} else {
		runServer()
	}
}

func runServer() {
	// Create a new echo instance
	e := echo.New()

	// Add a /courses endpoint
	e.GET("/courses", courses)

	// Start the server and throw
	// a fatal log if it fails
	e.Logger.Fatal(e.Start(":8000"))
}

func courses(c echo.Context) error {
	// Set the content response header to JSON
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	// Read courses.json, or throw a 500 error
	// if the file is not found
	data, err := ioutil.ReadFile("courses.json")
	if err != nil {
		errMessage := `{"error": "Could not load courses"}`
		return c.String(http.StatusInternalServerError, errMessage)
	}

	// Return the JSON in the response
	return c.String(http.StatusOK, string(data))
}
