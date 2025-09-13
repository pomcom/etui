package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	dataFile = "tasks.jsonl"
	appName = "etui"
)

func getDataFilePath() (string, error) {
	var dataDir string
	
	// Follow XDG Base Directory specification
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		dataDir = filepath.Join(xdgDataHome, appName)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(home, ".local", "share", "eisenhower-tui")
	}
	
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return "", err
	}
	
	return filepath.Join(dataDir, dataFile), nil
}

func (tm *TaskManager) SaveTasks() error {
	filePath, err := getDataFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, task := range tm.tasks {
		data, err := json.Marshal(task)
		if err != nil {
			return err
		}
		_, err = writer.WriteString(string(data) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (tm *TaskManager) LoadTasks() error {
	filePath, err := getDataFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	tm.tasks = make([]Task, 0)
	tm.nextID = 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var task Task
		err := json.Unmarshal([]byte(line), &task)
		if err != nil {
			continue
		}

		tm.tasks = append(tm.tasks, task)
		if task.ID >= tm.nextID {
			tm.nextID = task.ID + 1
		}
	}

	return scanner.Err()
}
