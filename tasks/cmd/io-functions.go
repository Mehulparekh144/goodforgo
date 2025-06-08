package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func getTasks(content []byte) []string {
	tasks := strings.Split(string(content), "\n")
	return tasks
}

func getTask(content []byte, ID int) (int, error) {
	tasks := getTasks(content)
	for i, task := range tasks {
		taskID := strings.Split(task, "\t")[0]

		intTaskId, err := strconv.Atoi(taskID)
		if err != nil {
			continue
		}

		if intTaskId == ID {
			return i, nil
		}
	}
	return -1, errors.New("task not found")
}

func DisplayTasks(content []byte) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(writer, "ID\tTime\tTask\tStatus")
	fmt.Fprintln(writer, "----\t----\t----\t----")
	fmt.Fprintln(writer, string(content))
	writer.Flush()
}

func AddTask(content []byte, file *os.File) {
	taskAdded := false

	lastTaskIDInt := 0
	tasks := getTasks(content)
	if len(tasks) > 0 {
		lastTask := tasks[0]

		if strings.TrimSpace(lastTask) != "" {
			lastTaskID := strings.Split(lastTask, "\t")[0]
			integer, err := strconv.Atoi(lastTaskID)
			if err != nil {
				fmt.Println("💥 Oops! Error converting last task ID to int:", err)
				os.Exit(1)
			}
			lastTaskIDInt = integer
		}
	}

	for !taskAdded {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n✍️  Enter the task to add \n (Press Enter to finish or type 'exit' to quit):")
		task, _ := reader.ReadString('\n')
		task = task[:len(task)-1]

		if strings.TrimSpace(task) == "exit" {
			fmt.Println("👋 Goodbye! See you next time!")
			os.Exit(0)
		}

		if strings.TrimSpace(task) == "" {
			fmt.Println("🤔 Hmm... looks like you forgot to enter a task!")
			fmt.Println("💡 Try again with a task description!")
		} else {
			now := time.Now()
			task = strconv.Itoa(lastTaskIDInt+1) + "\t" + now.Format(TIME_FORMAT) + "\t" + task + "\t" + "[ ]"

			newContent := []byte(strings.TrimSpace(task) + "\n")
			newContent = append(newContent, content...)
			_, err := file.WriteAt(newContent, 0)
			if err != nil {
				fmt.Println("💥 Oops! Error writing task to file: ", err)
			} else {
				fmt.Println("🎉 Task added successfully! You're getting things done! 💪")
				taskAdded = true
			}
		}
	}
}

func RemoveTask(content []byte, file *os.File, taskID string) {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		fmt.Println("❌ Error: Invalid task ID format")
		return
	}

	taskIndex, err := getTask(content, id)
	if err != nil {
		fmt.Println("❌ Oops! Task not found in our list!")
		return
	}

	tasks := getTasks(content)
	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
	newContent := []byte(strings.Join(tasks, "\n"))

	// Truncate the file before writing
	if err := file.Truncate(0); err != nil {
		fmt.Println("💥 Oops! Error truncating file:", err)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("💥 Oops! Error seeking to start of file:", err)
		return
	}

	_, err = file.Write(newContent)
	if err != nil {
		fmt.Println("💥 Oops! Error writing to file:", err)
		return
	}
	fmt.Println("🎉 Task removed successfully! One less thing to worry about!")
}

func CompleteTask(content []byte, file *os.File, taskID string) {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		fmt.Println("❌ Oops! That's not a valid task ID!")
		return
	}

	taskIndex, err := getTask(content, id)
	if err != nil {
		fmt.Println("❌ Oops! Task not found in our list!")
		return
	}

	tasks := getTasks(content)
	task := tasks[taskIndex]

	// Check if task is already completed
	if strings.Contains(task, "[X]") {
		fmt.Println("ℹ️ This task is already completed!")
		return
	}

	newTask := strings.Replace(task, "[ ]", "[X]", 1)
	tasks[taskIndex] = newTask

	newContent := []byte(strings.Join(tasks, "\n"))
	if len(newContent) > 0 {
		newContent = append(newContent, '\n')
	}

	// Truncate the file before writing
	if err := file.Truncate(0); err != nil {
		fmt.Println("💥 Oops! Error truncating file:", err)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("💥 Oops! Error seeking to start of file:", err)
		return
	}

	_, err = file.Write(newContent)
	if err != nil {
		fmt.Println("💥 Oops! Error writing to file:", err)
		return
	}
	fmt.Println("🎉 Task completed! You're on fire! 🔥")
}
