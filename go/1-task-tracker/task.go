package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type State int

const (
	Todo = iota
	InProgress
	Done
)

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	State State  `json:"state"`
}

type TaskManager struct {
	Tasks    []Task
	fileName string
}

// Custom String method for Task
func (t Task) String() string {
	return fmt.Sprintf("{Id: %d, Title: %s, State: %s}", t.Id, t.Title, t.State.String())
}

// Get the name of the constant as a string
func (s State) String() string {
	return [...]string{"To do", "In Progress", "Done"}[s]
}

// Get the name in snake case
func (s State) getName() string {
	name := s.String()
	if name == "" {
		return ""
	}

	var result strings.Builder
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: []Task{},
	}
}

// LoadTasks loads tasks from a JSON file
func (tm *TaskManager) LoadTasks(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("error when opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	// Check if the JSON is empty
	if decoder.More() {
		if err := decoder.Decode(&tm.Tasks); err != nil {
			return fmt.Errorf("error during reading data: %v", err)
		}
	} else {
		tm.Tasks = []Task{}
	}

	tm.fileName = filename
	return nil
}

func (tm *TaskManager) addTask(title string) *Task {
	newTask := Task{
		Id:    len(tm.Tasks) + 1,
		Title: title,
		State: Todo,
	}
	tm.Tasks = append(tm.Tasks, newTask)
	return &newTask
}

func (tm *TaskManager) updateTitle(id int, newTitle string) error {
	for i := range tm.Tasks {
		if tm.Tasks[i].Id == id {
			tm.Tasks[i].Title = newTitle
			return nil
		}
	}

	return fmt.Errorf("task with id %d not found", id)
}

func (tm *TaskManager) updateState(id int, state State) error {
	for i := range tm.Tasks {
		if tm.Tasks[i].Id == id {
			tm.Tasks[i].State = state
			return nil
		}
	}

	return fmt.Errorf("task with id %d not found", id)
}

func (tm *TaskManager) deleteTask(id int) error {
	for i := range tm.Tasks {
		if tm.Tasks[i].Id == id {
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("could not delete task with id %d", id)
}

func (tm *TaskManager) SaveTasks() error {
	file, err := os.OpenFile(tm.fileName, os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("error when opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tm.Tasks); err != nil {
		return fmt.Errorf("error during saving data: %v", err)
	}

	return nil
}
