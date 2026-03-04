package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Task      string `json:"task"`
	Status    string `json:"status"`
	CreatedAt string `json:"created at"`
}
type Command struct {
	command string
	run     func()
}

func change(file *os.File, tasks []Task) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("ошибка енкодинга файла: %v", err)
	}
	return nil
}
func newTask(tasks []Task, scanner *bufio.Scanner) []Task {
	var (
		taskName   string
		taskStatus string
	)
	fmt.Print("введи содержание таска: ")
	if scanner.Scan() {
		taskName = scanner.Text()
	} else {
		fmt.Printf("ошибка сканера")
		return tasks
	}
	fmt.Print("введи статус таска: ")
	if scanner.Scan() {
		taskStatus = scanner.Text()
	} else {
		fmt.Printf("ошибка сканера")
		return tasks
	}
	tasks = append(tasks, Task{taskName, taskStatus, time.Now().Format("2006-01-02 15:04")})
	return tasks
}
func deleteTask(tasks []Task, scanner *bufio.Scanner) []Task {
	checkTasks(tasks)
	var line int
	var err error
	fmt.Print("введи номер таска для удаления: ")
	if scanner.Scan() {
		line, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf("ошибка чтения номера таска (возможно вы ввели не числоа: %v)", err)
			return tasks
		}
	} else {
		fmt.Printf("ошибка сканера")
	}
	if line >= len(tasks) || line < 0 {
		fmt.Printf("ошибка: нет такой строки")
		return tasks
	}
	tasks = append(tasks[:line], tasks[line+1:]...)
	return tasks
}
func retype(tasks []Task, scanner *bufio.Scanner) []Task {
	var command int
	var err error
	var line int
	commands := []func(){
		func() {
			tasks, err = retypeTask(tasks, &line, scanner)
			if err != nil {
				fmt.Println(err)
			}
		},
		func() {
			tasks, err = retypeStatus(tasks, &line, scanner)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	checkTasks(tasks)
	fmt.Printf("введите номер таска который вы хотите изменить: ")
	if scanner.Scan() {
		line, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf("ошибка чтения строки (возможно вы ввели не число): %v", err)
			return tasks
		}
		if line >= len(tasks) || line < 0 {
			fmt.Printf("ошибка:  нет такой строки")
			return tasks
		}
	}
	fmt.Println("введите что вы хотите изменить:\n0. задача	\n1. статус")
	if scanner.Scan() {
		command, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Printf("ошибка чтения строки (возможно вы ввели не число): %v", err)
			return tasks
		}
	} else {
		fmt.Printf("ошибка сканера")
		return tasks
	}
	if command >= len(commands) || command < 0 {
		fmt.Println("нет такой команды")
		return tasks
	}
	commands[command]()
	return tasks
}
func retypeTask(tasks []Task, line *int, scanner *bufio.Scanner) ([]Task, error) {
	fmt.Printf("введите новое содержание задания: ")
	if scanner.Scan() {
		tasks[*line].Task = scanner.Text()
		return tasks, nil
	} else {
		return tasks, fmt.Errorf("ошибка сканера")
	}
}
func retypeStatus(tasks []Task, line *int, scanner *bufio.Scanner) ([]Task, error) {
	fmt.Printf("введите новое содержание статуса: ")
	if scanner.Scan() {
		tasks[*line].Status = scanner.Text()
		return tasks, nil
	} else {
		return tasks, fmt.Errorf("ошибка сканера")
	}
}
func checkTasks(tasks []Task) {
	for i, v := range tasks {
		fmt.Printf("%d. %s, статус: %v, дата создания: %s\n", i, v.Task, v.Status, v.CreatedAt)
	}
}
func main() {
	exit := false
	scanner := bufio.NewScanner(os.Stdin)
	var tasks []Task
	if _, err := os.Stat("tasks.json"); os.IsNotExist(err) {
		os.WriteFile("tasks.json", []byte("[]"), 0666)
	}
	jsonFile, err := os.ReadFile("tasks.json")
	if err != nil {
		fmt.Printf("ошибка чтения файла: %v", err)
		return
	}
	err = json.Unmarshal(jsonFile, &tasks)
	if err != nil {
		fmt.Printf("ошибка унмаршалинга: %v", err)
		return
	}
	var command int
	commandList := []Command{
		{"добавить задачу", func() { tasks = newTask(tasks, scanner) }},
		{"удалить задачу", func() { tasks = deleteTask(tasks, scanner) }},
		{"посмотреть задачи", func() { checkTasks(tasks) }},
		{"изменить задачу", func() { tasks = retype(tasks, scanner) }},
		{"выйти", func() {
			fmt.Println("до свидания!")
			exit = true
		}},
	}
	for !exit {
		fmt.Print("------------доступные команды------------")
		for i, v := range commandList {
			if i%2 == 0 {
				fmt.Println("")
			}
			fmt.Printf("%d. %s	", i, v.command)
		}
		fmt.Print("\nвведи номер команды: ")
		if scanner.Scan() {
			command, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Printf("ошибка чтения номера команды (возможно вы ввели не число): %v", err)
				continue
			}
		} else {
			fmt.Printf("ошибка сканера")
			return
		}
		if command >= len(commandList) || command < 0 {
			fmt.Printf("ошибка: нет такой команды")
			continue
		}
		commandList[command].run()
		file, err := os.Create("tasks.json")
		if err != nil {
			fmt.Printf("ошибка открытия файла: %v", err)
			continue
		}
		if err = change(file, tasks); err != nil {
			fmt.Println(err)
		}
		file.Close()
	}
}
