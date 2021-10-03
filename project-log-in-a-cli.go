package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ProjectLog struct {
	Day      string
	TaskList []string
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Adds an entry to the project log.
func saveToLog(filename string, entry *ProjectLog) {

	jsonData := []ProjectLog{}
	filteredJsonData := []ProjectLog{}
	dateAlreadyExists := false
	fileRead, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(fileRead, &jsonData)

	for _, object := range jsonData {
		if object.Day == fmt.Sprint(time.Now().Day())+
			"/"+fmt.Sprint(int(time.Now().Month()))+
			"/"+fmt.Sprint(time.Now().Year()) {
			object = ProjectLog{}
			object.Day = entry.Day
			object.TaskList = entry.TaskList
			filteredJsonData = append(filteredJsonData, object)
			dateAlreadyExists = true
			break
		}
		filteredJsonData = append(filteredJsonData, object)
	}

	if !dateAlreadyExists {
		filteredJsonData = append(filteredJsonData, *entry)
	}
	fileWrite, _ := json.MarshalIndent(filteredJsonData, "", " ")

	_ = ioutil.WriteFile(filename, fileWrite, 0644)

	// fmt.Println(string(fileWrite))
}

func dayStart(filename string) {

	dayEntry := &ProjectLog{
		Day: fmt.Sprint(time.Now().Day()) +
			"/" + fmt.Sprint(int(time.Now().Month())) +
			"/" + fmt.Sprint(time.Now().Year()),
		TaskList: []string{},
	}
	saveToLog(filename, dayEntry)
}

func dayEnd() {
	fmt.Println("day end test!")
}

func taskAdd(filename string, args []string) {

	dataTable := []map[string]interface{}{}
	file, err := ioutil.ReadFile(filename)
	// fileStr := string(file)
	// fmt.Println(fileStr)
	// fmt.Println(filename)
	// fmt.Println(args)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &dataTable)

	tasksAndTime := []string{}
	taskSet := ""
	for _, field := range args {
		if strings.HasPrefix(field, "task:") {
			taskSet = taskSet + strings.Split(field, ":")[1] + ":"
		} else if strings.HasPrefix(field, "time:") {
			taskSet = taskSet + strings.Split(field, ":")[1]
		}

		if len(strings.Split(taskSet, ":")[1]) > 0 {
			tasksAndTime = append(tasksAndTime, taskSet)
			taskSet = ""
		}
	}
	for _, data := range dataTable {
		if fmt.Sprint(time.Now().Day())+
			"/"+fmt.Sprint(int(time.Now().Month()))+
			"/"+fmt.Sprint(time.Now().Year()) == data["Day"] {
			fmt.Println(data["Day"])
			dayEntry := &ProjectLog{
				Day:      fmt.Sprint(data["Day"]),
				TaskList: tasksAndTime,
			}
			saveToLog(filename, dayEntry)
			break
		}
	}
}

func main() {
	args := os.Args

	switch args[1] {
	case "day":
		if args[2] == "start" {
			dayStart("log/project_log.json")
		} else if args[2] == "end" {
			dayEnd()
		}
	case "task":
		if args[2] == "add" {
			taskAdd("log/project_log.json", args[3:])
		}
	}
}
