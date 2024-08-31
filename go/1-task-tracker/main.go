package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	command, param, param2, err := getCommandsAndParams(os.Args)
	if err != nil {
		fmt.Printf("Error : %s", err)
		return
	}

	tm := NewTaskManager()
	if err := tm.LoadTasks("tasks.json"); err != nil {
		log.Fatal(err)
	}

	result := handleCommands(tm, command, param, param2)
	if err := tm.saveTasks(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func handleCommands(tm *TaskManager, command, param, param2 string) string {
	switch command {
	case "list":
		tasks := tasksToString(tm.Tasks, param)
		return tasks
	case "add":
		newTask := tm.addTask(param)
		return fmt.Sprintf("Task has been added successfully. (ID: %d)", newTask.Id)
	case "update":
		taskId, _ := strconv.Atoi(param)
		err := tm.updateTitle(taskId, param2)
		if err != nil {
			return fmt.Sprintf("Could not update the task, [error: %s]", err)
		}
		return fmt.Sprintf("Task has been updated successfully.")
	case "delete":
		taskId, _ := strconv.Atoi(param)
		err := tm.deleteTask(taskId)
		if err != nil {
			return fmt.Sprintf("Could not delete the task, [error: %s]", err)
		}
		return fmt.Sprintf("Task has been deleted successfully.")
	case "mark-in-progress":
		taskId, _ := strconv.Atoi(param)
		err := tm.updateState(taskId, InProgress)
		if err != nil {
			return fmt.Sprintf("Could not update state of the task, [error: %s]", err)
		}
	case "mark-done":
		taskId, _ := strconv.Atoi(param)
		err := tm.updateState(taskId, Done)
		if err != nil {
			return fmt.Sprintf("Could not update state of the task, [error: %s]", err)
		}
	default:
		return fmt.Sprintf("Command %s is not supported", command)
	}

	return fmt.Sprintf("")
}

func getCommandsAndParams(input []string) (c, p, p2 string, err error) {
	if len(input) < 2 {
		return "", "", "", fmt.Errorf("command is required")
	}

	// Skip the program name (at index 0)
	args := input[1:]

	if args[0] == "add" && len(args) > 1 {
		restOfArgs := strings.Join(args[1:], " ")
		args = []string{args[0], restOfArgs}
	}

	switch len(args) {
	case 1:
		return args[0], "", "", nil
	case 2:
		return args[0], args[1], "", nil
	case 3:
		return args[0], args[1], args[2], nil
	default:
		// If there are more than 3 arguments, combine the rest into the param
		return args[0], args[1], strings.Join(args[2:], " "), nil
	}
}

// Custom function to convert []Task to string
func tasksToString(tasks []Task, state string) string {
	var sb strings.Builder
	sb.WriteString("[\n")

	first := true

	for _, task := range tasks {
		if len(state) == 0 || (len(state) > 0 && task.State.getName() == state) {
			if !first {
				sb.WriteString(",\n")
			}
			sb.WriteString(" ")
			sb.WriteString(task.String())
			first = false
		}
	}

	sb.WriteString("\n]")
	return sb.String()
}
