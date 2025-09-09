package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	modeMatrix mode = iota
	modeAdd
	modeEdit
	modeDelete
)

type clearMessageMsg struct{}

func clearMessageAfter(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return clearMessageMsg{}
	})
}

func getQuadrantTip(quad Quadrant) string {
	switch quad {
	case UrgentImportant:
		return "DO FIRST"
	case NotUrgentImportant:
		return "SCHEDULE"
	case UrgentNotImportant:
		return "DELEGATE"
	case NotUrgentNotImportant:
		return "DON'T DO"
	default:
		return ""
	}
}

type model struct {
	taskManager       *TaskManager
	mode              mode
	selectedQuad      Quadrant
	selectedTaskIndex int
	selectedTaskID    int
	input             string
	editingField      int
	width             int
	height            int
	message           string
	showTips          bool

	tempTitle       string
	tempDescription string
	tempQuadrant    Quadrant
}

func initialModel() model {
	tm := NewTaskManager()
	tm.LoadTasks()

	return model{
		taskManager:  tm,
		mode:         modeMatrix,
		selectedQuad: UrgentImportant,
		width:        80,
		height:       24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case clearMessageMsg:
		m.message = ""
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeMatrix:
			return m.updateMatrix(msg)
		case modeAdd:
			return m.updateAdd(msg)
		case modeEdit:
			return m.updateEdit(msg)
		case modeDelete:
			return m.updateDelete(msg)
		}
	}

	return m, nil
}

func (m model) updateMatrix(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.taskManager.SaveTasks()
		return m, tea.Quit

	case "t":
		m.taskManager.ToggleContext()
		m.selectedTaskIndex = 0
		m.message = ""
		
	case "?", "f1":
		m.showTips = !m.showTips
		if m.showTips {
			m.message = "Tips enabled"
		} else {
			m.message = "Tips disabled"
		}
		return m, clearMessageAfter(1*time.Second)

	case "h", "left":
		if m.selectedQuad == UrgentImportant {
			m.selectedQuad = NotUrgentImportant
		} else if m.selectedQuad == UrgentNotImportant {
			m.selectedQuad = NotUrgentNotImportant
		} else if m.selectedQuad == NotUrgentImportant {
			m.selectedQuad = UrgentImportant
		} else if m.selectedQuad == NotUrgentNotImportant {
			m.selectedQuad = UrgentNotImportant
		}
		m.selectedTaskIndex = 0

	case "l", "right":
		if m.selectedQuad == UrgentImportant {
			m.selectedQuad = NotUrgentImportant
		} else if m.selectedQuad == NotUrgentImportant {
			m.selectedQuad = UrgentImportant
		} else if m.selectedQuad == UrgentNotImportant {
			m.selectedQuad = NotUrgentNotImportant
		} else if m.selectedQuad == NotUrgentNotImportant {
			m.selectedQuad = UrgentNotImportant
		}
		m.selectedTaskIndex = 0

	case "j", "down":
		tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
		if len(tasks) > 0 && m.selectedTaskIndex < len(tasks)-1 {
			m.selectedTaskIndex++
		} else {
			if m.selectedQuad == UrgentImportant {
				m.selectedQuad = UrgentNotImportant
			} else if m.selectedQuad == NotUrgentImportant {
				m.selectedQuad = NotUrgentNotImportant
			} else if m.selectedQuad == UrgentNotImportant {
				m.selectedQuad = UrgentImportant
			} else if m.selectedQuad == NotUrgentNotImportant {
				m.selectedQuad = NotUrgentImportant
			}
			m.selectedTaskIndex = 0
		}

	case "k", "up":
		if m.selectedTaskIndex > 0 {
			m.selectedTaskIndex--
		} else {
			if m.selectedQuad == UrgentImportant {
				m.selectedQuad = UrgentNotImportant
			} else if m.selectedQuad == NotUrgentImportant {
				m.selectedQuad = NotUrgentNotImportant
			} else if m.selectedQuad == UrgentNotImportant {
				m.selectedQuad = UrgentImportant
			} else if m.selectedQuad == NotUrgentNotImportant {
				m.selectedQuad = NotUrgentImportant
			}
			tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
			if len(tasks) > 0 {
				m.selectedTaskIndex = len(tasks) - 1
			} else {
				m.selectedTaskIndex = 0
			}
		}

	case "a":
		m.mode = modeAdd
		m.tempQuadrant = m.selectedQuad
		m.tempTitle = ""
		m.tempDescription = ""
		m.editingField = 0
		m.input = ""
		m.message = ""

	case "d":
		tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
		if len(tasks) > m.selectedTaskIndex {
			m.mode = modeDelete
			m.selectedTaskID = tasks[m.selectedTaskIndex].ID
			m.message = ""
		}

	case "e":
		tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
		if len(tasks) > m.selectedTaskIndex {
			task := tasks[m.selectedTaskIndex]
			m.mode = modeEdit
			m.selectedTaskID = task.ID
			m.tempTitle = task.Title
			m.tempDescription = task.Description
			m.tempQuadrant = task.Quadrant
			m.editingField = 0
			m.input = task.Title
			m.message = ""
		}

	case " ", "enter":
		tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
		if len(tasks) > m.selectedTaskIndex {
			task := tasks[m.selectedTaskIndex]
			m.taskManager.ToggleTaskCompletion(task.ID)
			m.taskManager.SaveTasks()
			if task.Completed {
				m.message = "Task marked as incomplete! 🔄"
			} else {
				m.message = "Task completed! ✅"
			}
			return m, clearMessageAfter(3*time.Second)
		}
	}

	return m, nil
}

func (m model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeMatrix
		m.message = "Add cancelled"

	case "enter":
		if m.editingField == 0 {
			m.tempTitle = m.input
			m.editingField = 1
			m.input = ""
		} else {
			m.tempDescription = m.input
			m.taskManager.AddTask(m.tempTitle, m.tempDescription, m.tempQuadrant, m.taskManager.GetCurrentContext())
			m.taskManager.SaveTasks()
			m.mode = modeMatrix
			m.message = "Task added successfully"
			return m, clearMessageAfter(3*time.Second)
		}

	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}

	return m, nil
}

func (m model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeMatrix
		m.message = "Edit cancelled"

	case "enter":
		if m.editingField == 0 {
			m.tempTitle = m.input
			m.editingField = 1
			m.input = m.tempDescription
		} else {
			m.tempDescription = m.input
			m.taskManager.UpdateTask(m.selectedTaskID, m.tempTitle, m.tempDescription, m.tempQuadrant)
			m.taskManager.SaveTasks()
			m.mode = modeMatrix
			m.message = "Task updated successfully"
			return m, clearMessageAfter(3*time.Second)
		}

	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}

	return m, nil
}

func (m model) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		m.mode = modeMatrix
		m.message = "Delete cancelled"

	case "y":
		m.taskManager.DeleteTask(m.selectedTaskID)
		m.taskManager.SaveTasks()

		tasks := m.taskManager.GetTasksByQuadrant(m.selectedQuad)
		if m.selectedTaskIndex >= len(tasks) && len(tasks) > 0 {
			m.selectedTaskIndex = len(tasks) - 1
		} else if len(tasks) == 0 {
			m.selectedTaskIndex = 0
		}

		m.mode = modeMatrix
		m.message = "Task deleted successfully"
		return m, clearMessageAfter(3*time.Second)
	}

	return m, nil
}