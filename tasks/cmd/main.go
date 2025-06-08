package main

import (
	"fmt"
	"io"
	"os"
)

const FILE_PATH = "/home/mparekh/Desktop/Coding/go-proj/tasks/files/tasks.txt"
const ASCII_PATH = "/home/mparekh/Desktop/Coding/go-proj/tasks/files/ascii.txt"

func main() {

	ascii, err := os.ReadFile(ASCII_PATH)
	if err != nil {
		fmt.Println("Error reading ASCII file:", err)
		os.Exit(1)
	}

	fmt.Print(string(ascii) + "\n")

	args := os.Args
	if len(args) > 1 {

		file, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error reading tasks file:", err)
			os.Exit(1)
		}

		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading tasks file:", err)
			os.Exit(1)
		}

		switch args[1] {

		case ".":
			DisplayTasks(content)
		case "+":
			AddTask(content, file)

		case "-":
			if len(args) < 3 {
				fmt.Println("Please provide a task ID to remove")
				os.Exit(1)
			}
			RemoveTask(content, file, args[2])

		case "!":
			if len(args) < 3 {
				fmt.Println("Please provide a task ID to complete")
				os.Exit(1)
			}
			CompleteTask(content, file, args[2])
		}

	} else {
		fmt.Print("\n🎯 Welcome To Tasks 🎯\n" +
			"Please provide args to run, e.g. 'go run main.go <args>'\n" +
			"Available args: \n" +
			"  - . 📋 : List all tasks\n" +
			"  - + ➕ : Add a new task\n" +
			"  - - ❌ <task_id> : Remove a task by ID\n" +
			"  - ! ✅ <task_id> : Mark a task as complete\n")
	}

}
