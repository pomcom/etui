package main

import (
	"time"
)

type Context string

const (
	ContextWork    Context = "work"
	ContextPrivate Context = "private"
)

type Quadrant int

const (
	UrgentImportant Quadrant = iota
	NotUrgentImportant
	UrgentNotImportant
	NotUrgentNotImportant
)

func (q Quadrant) String() string {
	switch q {
	case UrgentImportant:
		return "Urgent & Important"
	case NotUrgentImportant:
		return "Not Urgent & Important"
	case UrgentNotImportant:
		return "Urgent & Not Important"
	case NotUrgentNotImportant:
		return "Not Urgent & Not Important"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Quadrant    Quadrant  `json:"quadrant"`
	Context     Context   `json:"context"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskManager struct {
	tasks       []Task
	nextID      int
	currentView Context
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:       make([]Task, 0),
		nextID:      1,
		currentView: ContextWork,
	}
}

func (tm *TaskManager) AddTask(title, description string, quadrant Quadrant, context Context) *Task {
	task := Task{
		ID:          tm.nextID,
		Title:       title,
		Description: description,
		Quadrant:    quadrant,
		Context:     context,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	return &task
}

func (tm *TaskManager) GetTasks() []Task {
	var filtered []Task
	for _, task := range tm.tasks {
		if task.Context == tm.currentView {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func (tm *TaskManager) GetTasksByQuadrant(quadrant Quadrant) []Task {
	var filtered []Task
	for _, task := range tm.tasks {
		if task.Context == tm.currentView && task.Quadrant == quadrant {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func (tm *TaskManager) UpdateTask(id int, title, description string, quadrant Quadrant) bool {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks[i].Title = title
			tm.tasks[i].Description = description
			tm.tasks[i].Quadrant = quadrant
			tm.tasks[i].UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (tm *TaskManager) DeleteTask(id int) bool {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (tm *TaskManager) ToggleTaskCompletion(id int) bool {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks[i].Completed = !tm.tasks[i].Completed
			tm.tasks[i].UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (tm *TaskManager) ToggleContext() {
	if tm.currentView == ContextWork {
		tm.currentView = ContextPrivate
	} else {
		tm.currentView = ContextWork
	}
}

func (tm *TaskManager) GetCurrentContext() Context {
	return tm.currentView
}