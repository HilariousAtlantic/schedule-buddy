package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"

	"github.com/lib/pq"
)

const (
	databasePath = "user=schedule_buddy dbname=schedule_buddy sslmode=disable"
)

const (
	createCoursesTable = `
	CREATE TABLE courses (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		subject TEXT NOT NULL,
		number TEXT NOT NULL,
		credits TEXT NOT NULL
	);
	`
	insertCourse = `
	INSERT INTO courses (name, subject, number, credits) VALUES (?, ?, ?, ?)
	`

	selectCourses = `
	SELECT id, name, subject, number, credits FROM courses
	`
)

type DB struct {
	db *sql.DB
}

var dbContext = new(DB)

func (d *DB) open() *sql.DB {
	if d.db == nil {
		var err error
		d.db, err = sql.Open("postgres", databasePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	return d.db
}

func createDatabase() {
	fmt.Println("Creating database...")

	output, err := exec.Command("createdb", "schedule_buddy", "-U", "schedule_buddy").CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		log.Fatal(err)
	} else {
		fmt.Println("Database created")
	}

	db := dbContext.open()
	defer db.Close()

	_, err = db.Exec(createCoursesTable)
	if err != nil {
		log.Fatalf("%q: %s\n", err, createCoursesTable)
	}
}

func batchInsertCourses(courses []*Course) {
	existingCourseNames := map[string]bool{}
	db := dbContext.open()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(pq.CopyIn("courses", "name", "subject", "number", "credits"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, course := range courses {
		if existingCourseNames[course.Name] {
			continue
		} else {
			existingCourseNames[course.Name] = true
		}
		_, err = stmt.Exec(course.Name, course.Subject, course.Number, course.Credits)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func getCoursesFromDB() []*Course {
	var courses []*Course
	db := dbContext.open()
	rows, err := db.Query(selectCourses)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		course := &Course{}
		err = rows.Scan(&course.ID, &course.Name, &course.Subject, &course.Number, &course.Credits)
		if err != nil {
			log.Fatal(err)
		}
		courses = append(courses, course)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return courses
}

func deleteDatabase() {
	fmt.Println("Deleting database...")

	output, err := exec.Command("dropdb", "schedule_buddy", "-U", "schedule_buddy").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error(), string(output))
	} else {
		fmt.Println("Database deleted")
	}
}
