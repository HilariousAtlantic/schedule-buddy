package main

import (
	"database/sql"
	//	"fmt"
	"math"
	"strconv"
	"strings"
)

// a helper function that is used to create the necessary
// variables to call findGoodSchedulesRecursive
func findGoodSchedules(ids string) []Schedule {

	goodSchedules := make([]Schedule, 0)

	//courses cannot be a pointer due to recursion, so convert
	coursesPointer := getCourseTree(ids)
	setCoursesGPA(coursesPointer)

	courses := make([]Course, 0)
	for _, coursePointer := range coursesPointer {
		courses = append(courses, *coursePointer)
	}

	selectedSections := make([]Section, 0)
	findGoodSchedulesRecursive(courses, selectedSections, &goodSchedules)

	//fmt.Println(goodSchedules)
	//fmt.Println("good schedules")

	return goodSchedules

}

/*recursive function to find all good schedules.
  works by checking every combination until it finds an invalid time, then returning and checking
  the next section.
*/
func findGoodSchedulesRecursive(courses []Course, selectedSections []Section, goodSchedules *[]Schedule) {
	//base case; we found a good schedule. Append it and return.
	if len(selectedSections) == len(courses) {
		var sections []ScheduledCourse
		goodScheduleAvgGPA := 0.0
		divideBy := 0.0
		for _, selectedSection := range selectedSections {
			//	fmt.Println(selectedSection.AverageGPA)
			if !(selectedSection.AverageGPA == 0.0) {
				goodScheduleAvgGPA += (selectedSection.AverageGPA * selectedSection.Credits)
				divideBy += selectedSection.Credits
				//fmt.Printf("goodSchedGPA is now : %v, divideBy is now: %v", goodScheduleAvgGPA, divideBy)
			}
			var scheduledCourse = ScheduledCourse{
				CourseID:  selectedSection.CourseID,
				SectionID: selectedSection.ID,
			}
			sections = append(sections, scheduledCourse)
		}
		avgGPA := 0.0
		if divideBy != 0 {
			avgGPA = goodScheduleAvgGPA / divideBy
		}
		avgGPA = float64(int(avgGPA*100)) / 100

		var goodSchedule = Schedule{
			Sections:   sections,
			AverageGPA: (avgGPA),
		}

		//fmt.Printf("sched gpa is: %v", goodSchedule.AverageGPA)
		*goodSchedules = append(*goodSchedules, goodSchedule)
		return
	}
	//skip a course and continue
SKIPCOURSE:
	for _, course := range courses {
		for _, selectedSection := range selectedSections {
			//if we have a course in the selectedSections, we dont check the other sections of said course
			if selectedSection.CourseID == course.ID {
				continue SKIPCOURSE
			}
		}
	SKIPSECTION:
		for _, section := range course.Sections {
			//go through all selectedSections and make sure none overlap
			for _, selectedSection := range selectedSections {

				//if overlap, skip that section
				if doTimesOverlap(selectedSection, *section) {
					continue SKIPSECTION
				}
			}
			//if none overlap, section is good, add to selectedSections
			selectedSections = append(selectedSections, *section)
			findGoodSchedulesRecursive(courses, selectedSections, goodSchedules)
			selectedSections = selectedSections[:len(selectedSections)-1]
		}
		return
	}
	return
}

//returns a built out "tree" of courses.
//a tree means that the courses have sections and the sections have meets
func getCourseTree(ids string) []*Course {

	courses := getCoursesFromIDs(ids)
	sections := getSectionsFromCourses(courses)
	meets := getMeetsFromSections(sections)

	for _, section := range sections {
		if len(meets) == 0 {
			break
		}
		for _, meet := range meets {
			if meet.SectionID == section.ID {
				section.Meets = append(section.Meets, meet)
			}
		}
	}
	for _, course := range courses {
		for _, section := range sections {
			if section.CourseID == course.ID {
				course.Sections = append(course.Sections, section)
			}
		}
	}
	return courses
}

func setCoursesGPA(courses []*Course) {
	for _, course := range courses {
		for _, section := range course.Sections {
			//credits stuff
			creditsStr := course.Credits
			//check for credit range
			if strings.Index(creditsStr, "-") != -1 {
				temp := strings.Split(creditsStr, "-")
				creditsStr = temp[0]
			}
			credits, err := strconv.ParseFloat(creditsStr, 64)
			handleError(err)
			section.Credits = credits

			//GPA stuff
			meet := section.Meets[0]
			instructor := meet.Instructor
			index := strings.Index(instructor, ";")

			//there is two professors
			if index != -1 {
				instructors := strings.Split(instructor, ";")
				//instructor1 := instructors[0]
				//instructors2 := instructors[1]
				avgGPA := (getAvgGPA(instructors[0], *course) + getAvgGPA(instructors[1], *course)) / 2
				section.AverageGPA = avgGPA
			} else {
				//	fmt.Println(getAvgGPA(instructor, *course))
				section.AverageGPA = getAvgGPA(instructor, *course)
				//	fmt.Printf("setCourse avg gpa is: %v", section.AverageGPA)
			}
		}
	}
}

func getAvgGPA(instructor string, course Course) float64 {
	var avgGPA = 0.0
	var divideBy = 0.0
	instructor = strings.TrimSpace(instructor)
	//	fmt.Printf("instructor:%v.\n", instructor)
	db := dbContext.open()
	var rows *sql.Rows
	var err error
	query := "SELECT gpa FROM grades WHERE LOWER(instructor) LIKE LOWER('%" + instructor +
		"%') AND subject = '" + course.Subject +
		"' AND number = '" + course.Number + "';"
	//query := "SELECT gpa FROM grades WHERE LOWER(instructor) LIKE LOWER('%zmuda%')"
	rows, err = db.Query(query)
	handleError(err)
	defer rows.Close()
	for rows.Next() {
		divideBy++
		var gpa float64
		err = rows.Scan(&gpa)
		//		fmt.Printf("in loop gpa: %v", gpa)
		avgGPA += gpa
		handleError(err)
	}
	err = rows.Err()
	handleError(err)
	avgGPA = (avgGPA / divideBy)
	avgGPA = float64(int(avgGPA*100)) / 100
	if avgGPA < 0 || avgGPA > 4 || math.IsNaN(avgGPA) {
		avgGPA = 0.0
	}
	//	fmt.Printf("avg gpa is: %v", avgGPA)
	return avgGPA

}

//checks if two sections have any meet times that over lap,
// ignoring meets that have no overlapping days
func doTimesOverlap(a, b Section) bool {
	for _, meetA := range a.Meets {
		for _, meetB := range b.Meets {
			//fmt.Printf("comparing meet: %v to meet: %v", meetA, meetB)
			if !containsSameDay(meetA.Days, meetB.Days) {
				//fmt.Println("different day")
				continue
			} else if meetA.StartTime <= meetB.EndTime &&
				meetB.StartTime <= meetA.EndTime {
				//fmt.Println(meetB)
				//fmt.Println(meetA)
				//fmt.Println("no overlap")
				return true
			}
		}
	}
	return false
}

//returns true if the two strings share a similar character
//accounts for TBA and any upper/lower case issues
func containsSameDay(a, b string) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	if a == "tba" || b == "tba" {
		return false
	}
	m := make(map[string]bool)
	for _, c := range a {
		s := string(c)
		m[s] = true
	}
	for _, c := range b {
		s := string(c)
		val, _ := m[s]
		if val == true {
			return true
		}
	}
	return false
}
