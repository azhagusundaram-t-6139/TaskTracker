package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

var taskList []Task
var latestTaskID = 0

func main() {
	ReadFile()
	initId()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide valid inputs")
		os.Exit(0)
	}

	operation, err := ParseOps(args[1])
	if err != nil {
		os.Exit(0)
	}

	switch operation {
	case Add:
		if len(args) != 3 {
			fmt.Println("Expected number of input arguments are 2")
			os.Exit(1)
		}
		AddTask(args[2])
		WriteFile()
	case Update:
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Please provide integer as id:", err)
			os.Exit(1)
		}
		UpdateTask(id, args[3])
		WriteFile()
	case List:
		if len(args) == 2 {
			ListTasks("")
			return
		}
		ListTasks(Status(args[2]))
	case Delete:
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Please provide integer as id:", err)
			os.Exit(1)
		}
		DeleteTask(id)
		WriteFile()
	case MarkInProgress:
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Please provide integer as id:", err)
			os.Exit(1)
		}
		MakeInProgress(id)
		WriteFile()
	case MarkDone:
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Please provide integer as id")
			os.Exit(0)
		}
		MakeDone(id)
		WriteFile()
	default:
		fmt.Println("Invalid operation")
		os.Exit(0)
	}
}

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
)

type Ops string

const (
	Add            Ops = "add"
	Update         Ops = "update"
	List           Ops = "list"
	Delete         Ops = "delete"
	MarkInProgress Ops = "mark-in-progress"
	MarkDone       Ops = "mark-done"
)

func ParseOps(s string) (Ops, error) {
	operation := Ops(s)

	switch operation {
	case Add, Update, List, Delete, MarkDone, MarkInProgress:
		return operation, nil
	default:
		return "", fmt.Errorf("Please provide a valid operation")
	}
}

type Task struct {
	ID          int       `json: "id"`
	Description string    `json: "description"`
	Status      Status    `json: "status"`
	CreatedAt   time.Time `json: "createdAt"`
	UpdatedAt   time.Time `json: "updatedAt"`
}

func AddTask(description string) []Task {
	task := Task{
		ID:          latestTaskID,
		Description: description,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	taskList = append(taskList, task)
	latestTaskID++
	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	return taskList
}

func UpdateTask(id int, desc string) {
	taskListAdr := &taskList
	for index, task := range *taskListAdr {
		if task.ID == id {
			(*taskListAdr)[index].Description = desc
		}
	}
}

func DeleteTask(id int) {
	taskList = slices.DeleteFunc(taskList, func(t Task) bool {
		return t.ID == id
	})
}

func ListTasks(status Status) []Task {
	switch status {
	case Todo:
		var openTasks []Task
		for _, task := range taskList {
			if task.Status == Todo {
				openTasks = append(openTasks, task)
			}
		}
		PrintResult(openTasks)
		return openTasks
	case Done:
		var doneTasks []Task
		for _, task := range taskList {
			if task.Status == Done {
				doneTasks = append(doneTasks, task)
			}
		}
		PrintResult(doneTasks)
		return doneTasks
	case InProgress:
		var inProgressTasks []Task
		for _, task := range taskList {
			if task.Status == InProgress {
				inProgressTasks = append(inProgressTasks, task)
			}
		}
		PrintResult(inProgressTasks)
		return inProgressTasks
	default:
		PrintResult(taskList)
		return taskList
	}
}

func MakeInProgress(id int) {
	if len(taskList) < id {
		fmt.Println("No task available in the list with id:", id)
		return
	}

	taskListAdr := &taskList
	for index, i := range *taskListAdr {
		if i.ID == id {
			(*taskListAdr)[index].Status = InProgress
		}
	}
}

func MakeDone(id int) {
	if len(taskList) < id {
		fmt.Println("No task available in the list with id:", id)
		return
	}

	taskListAdr := &taskList
	for index, i := range *taskListAdr {
		if i.ID == id {
			(*taskListAdr)[index].Status = Done
		}
	}
}

func WriteFile() {
	jsonData, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		fmt.Println("Unable to marshal the task list into json:", err)
		os.Exit(1)
	}

	err = os.WriteFile("todo.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
}

func ReadFile() {
	data, err := os.ReadFile("todo.json")
	if err != nil {
		fmt.Println("Error reading the todo json file; may be file not present")
		return
	}

	if len(data) == 0 {
		fmt.Println("File is empty. Starting with an empty task list.")
		return
	}

	err = json.Unmarshal(data, &taskList)
	if err != nil {
		fmt.Println("Unable to unmarshal the json data:", err)
		os.Exit(1)
	}
}

func PrintResult(taskList []Task) {
	jsonData, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		fmt.Println("Unable to marshal the task list into json:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func initId() {
	idSlice := []int{}
	if len(taskList) == 0 {
		latestTaskID = 0
		return
	}
	for _, t := range taskList {
		idSlice = append(idSlice, t.ID)
	}
	latestTaskID = slices.Max(idSlice) + 1
}
