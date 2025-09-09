# etui

**E**isenhower **T**erminal **U**ser **I**nterface - A minimal TUI for task management using the Eisenhower Matrix.

## Features

- **Eisenhower Matrix**: Organize tasks into 4 quadrants based on urgency and importance
- **Context Toggle**: Switch between work and private task contexts  
- **CRUD Operations**: Create, Read, Update, and Delete tasks
- **XDG Compliant**: Follows Unix standards for data storage
- **Persistent Storage**: Tasks saved in JSONL format
- **Keyboard Navigation**: Full keyboard-driven interface

## Installation

Install directly with Go:
```bash
go install github.com/pomcom/etui@latest
```

Or build from source:
```bash
git clone https://github.com/pomcom/etui.git
cd etui
go build -o etui .
```

## Usage

Run the application:
```bash
etui
```

### Controls

- **Navigation**: `←` `→` / `h` `l` to change quadrant, `↑` `↓` / `j` `k` to navigate tasks
- **Complete Task**: `Space` or `Enter` - Mark task as completed/incomplete
- **Add Task**: `a` - Add a new task to the selected quadrant
- **Edit Task**: `e` - Edit the selected task
- **Delete Task**: `d` - Delete the selected task
- **Toggle Context**: `t` - Switch between work and private contexts
- **Quit**: `q` or `Ctrl+C`

### Eisenhower Matrix Quadrants

1. **Urgent & Important** (Do First): Deadlines
2. **Not Urgent & Important** (Schedule): Planning, prevention, development
3. **Urgent & Not Important** (Delegate): Interruptions
4. **Not Urgent & Not Important** (Don't Do): Time wasters

## Data Storage

Tasks are automatically saved following XDG Base Directory specification:
- `$XDG_DATA_HOME/etui/tasks.jsonl` (if `XDG_DATA_HOME` is set)  
- `~/.local/share/etui/tasks.jsonl` (default)

Data is stored in JSONL (JSON Lines) format. Each line contains a JSON object representing a task:

```json
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "quadrant": 0,
  "context": "work",
  "completed": false,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
