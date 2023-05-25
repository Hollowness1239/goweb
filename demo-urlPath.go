package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor"`
}

var CourseList []Course

func init() {
	CourseJSON := `[
		{
			"id":1,
			"name":"python",
			"price":2590,
			"instructor": "thuu"
		},
		{
			"id":2,
			"name":"go",
			"price":5000,
			"instructor": "pita"
		},
		{
			"id":1,
			"name":"flutter",
			"price":10000,
			"instructor": "boss"
		}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	heightID := -1
	for _, course := range CourseList {
		if heightID < course.ID {
			heightID = course.ID
		}
	}
	return heightID + 1
}

func findID(ID int) (*Course, int) {
	for i, course := range CourseList {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	findCourse, index := findID(ID)
	fmt.Println("findCourse:", findCourse)
	fmt.Println("index:", index)
	if findCourse == nil {
		http.Error(w, fmt.Sprintf("no course with ID %d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(findCourse)
		fmt.Println("courseJSON:", courseJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPut:
		var updateCourse Course
		byteBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateCourse.ID == ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		findCourse = &updateCourse
		CourseList[index] = *findCourse
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON2, err := json.Marshal(CourseList)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON2)
	case http.MethodPost:
		var newCourse Course
		Bodybyte, err := ioutil.ReadAll(r.Body)
		fmt.Println("Bodybyte:", string(Bodybyte))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err2 := json.Unmarshal(Bodybyte, &newCourse)
		fmt.Println("newCourse", newCourse)
		if err2 != nil {
			fmt.Println("err2:", err2)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.HandleFunc("/course", coursesHandler)
	http.ListenAndServe(":5000", nil)
}
