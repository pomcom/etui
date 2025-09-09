package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.mode {
	case modeMatrix:
		return m.viewMatrix()
	case modeAdd:
		return m.viewAdd()
	case modeEdit:
		return m.viewEdit()
	case modeDelete:
		return m.viewDelete()
	default:
		return ""
	}
}

func (m model) viewMatrix() string {
	var s strings.Builder

	var titleColor string
	if m.taskManager.GetCurrentContext() == ContextWork {
		titleColor = "#7D56F4" // Purple for work
	} else {
		titleColor = "#04B575" // Green for private
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color(titleColor)).
		Padding(0, 1).
		Render(fmt.Sprintf("Eisenhower Matrix - %s Context", strings.Title(string(m.taskManager.GetCurrentContext()))))

	s.WriteString(title + "\n\n")

	quadWidth := (m.width - 6) / 2
	quadHeight := (m.height - 10) / 2

	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderQuadrant(UrgentImportant, quadWidth, quadHeight),
		" │ ",
		m.renderQuadrant(NotUrgentImportant, quadWidth, quadHeight),
	)

	separator := strings.Repeat("─", m.width-2)

	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderQuadrant(UrgentNotImportant, quadWidth, quadHeight),
		" │ ",
		m.renderQuadrant(NotUrgentNotImportant, quadWidth, quadHeight),
	)

	matrix := lipgloss.JoinVertical(lipgloss.Left, topRow, separator, bottomRow)
	s.WriteString(matrix)

	if m.message != "" {
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(m.message))
	}

	help := "\n\nControls: ← → / hl: change quadrant | ↑ ↓ / jk: navigate tasks | Space/Enter: complete | a: add | e: edit | d: delete | t: toggle context | ?: toggle tips | q: quit"
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(help))

	return s.String()
}

func (m model) renderQuadrant(quad Quadrant, width, height int) string {
	tasks := m.taskManager.GetTasksByQuadrant(quad)

	var style lipgloss.Style
	var borderStyle lipgloss.Border
	var activeColor, inactiveColor string

	// Priority-based warm colors
	var priorityColor string
	switch quad {
	case UrgentImportant:
		priorityColor = "#E07A3F" // Warm coral/amber
	case NotUrgentImportant:
		priorityColor = "#2A7F7C" // Deep teal
	case UrgentNotImportant:
		priorityColor = "#D4A574" // Warm gold
	case NotUrgentNotImportant:
		priorityColor = "#8B8680" // Soft gray
	}

	if m.taskManager.GetCurrentContext() == ContextWork {
		borderStyle = lipgloss.RoundedBorder()
		activeColor = priorityColor
		inactiveColor = "#3C3C3C"
	} else {
		borderStyle = lipgloss.DoubleBorder()
		activeColor = priorityColor
		inactiveColor = "#2C4A3D"
	}

	if quad == m.selectedQuad {
		style = lipgloss.NewStyle().
			Width(width).
			Height(height).
			Border(borderStyle).
			BorderForeground(lipgloss.Color(activeColor)).
			Padding(1)
	} else {
		style = lipgloss.NewStyle().
			Width(width).
			Height(height).
			Border(borderStyle).
			BorderForeground(lipgloss.Color(inactiveColor)).
			Padding(1)
	}

	var content strings.Builder
	header := quad.String()
	if m.showTips {
		tip := getQuadrantTip(quad)
		tipStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(priorityColor)).
			Italic(true)
		header = quad.String() + "\n" + tipStyle.Render(tip)
	}
	content.WriteString(lipgloss.NewStyle().Bold(true).Render(header) + "\n")

	maxTasks := height - 3
	if m.showTips {
		maxTasks = height - 4 // Account for extra tip line
	}
	for i, task := range tasks {
		if i >= maxTasks {
			content.WriteString("...")
			break
		}

		var taskStyle lipgloss.Style
		var prefix string
		var taskTitle string

		if task.Completed {
			if quad == m.selectedQuad && i == m.selectedTaskIndex {
				taskStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#000000")).
					Background(lipgloss.Color(activeColor)).
					Bold(true).
					Strikethrough(true)
				prefix = "✓"
			} else {
				taskStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#666666")).
					Strikethrough(true)
				prefix = "✓"
			}
			taskTitle = task.Title
		} else {
			if quad == m.selectedQuad && i == m.selectedTaskIndex {
				taskStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#000000")).
					Background(lipgloss.Color(activeColor)).
					Bold(true)
				prefix = "▶"
			} else {
				taskStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
				prefix = "•"
			}
			taskTitle = task.Title
		}

		taskLine := fmt.Sprintf("%s %s", prefix, taskTitle)
		content.WriteString(taskStyle.Render(taskLine) + "\n")

		if quad == m.selectedQuad && i == m.selectedTaskIndex && task.Description != "" {
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Italic(true)

			maxDescWidth := width - 4
			if len(task.Description) > maxDescWidth {
				truncated := task.Description[:maxDescWidth-3] + "..."
				content.WriteString(descStyle.Render("  "+truncated) + "\n")
			} else {
				content.WriteString(descStyle.Render("  "+task.Description) + "\n")
			}
		}
	}

	return style.Render(content.String())
}

func (m model) viewAdd() string {
	var s strings.Builder

	s.WriteString("Add New Task\n\n")

	if m.editingField == 0 {
		s.WriteString(fmt.Sprintf("Title: %s█", m.input))
	} else {
		s.WriteString(fmt.Sprintf("Title: %s", m.tempTitle))
		s.WriteString(fmt.Sprintf("\nDescription: %s█", m.input))
	}

	s.WriteString(fmt.Sprintf("\nQuadrant: %s", m.tempQuadrant.String()))
	s.WriteString("\n\nPress Enter to confirm, Esc to cancel")

	return s.String()
}

func (m model) viewEdit() string {
	var s strings.Builder

	s.WriteString("Edit Task\n\n")

	if m.editingField == 0 {
		s.WriteString(fmt.Sprintf("Title: %s█", m.input))
	} else {
		s.WriteString(fmt.Sprintf("Title: %s", m.tempTitle))
		s.WriteString(fmt.Sprintf("\nDescription: %s█", m.input))
	}

	s.WriteString(fmt.Sprintf("\nQuadrant: %s", m.tempQuadrant.String()))
	s.WriteString("\n\nPress Enter to confirm, Esc to cancel")

	return s.String()
}

func (m model) viewDelete() string {
	return fmt.Sprintf("Delete task with ID %d?\n\nPress 'y' to confirm, 'n' or Esc to cancel", m.selectedTaskID)
}