# Task Manager CLI

A simple command-line task manager written in Go that helps you keep track of your tasks.

## Features

- Add new tasks
- Remove tasks
- Mark tasks as complete
- Display tasks in a nicely formatted table
- Persistent storage of tasks

## Building from Source

### Prerequisites

- Go 1.16 or higher
- Git (optional, for cloning the repository)

### Build Instructions

#### Linux/macOS

```bash
# Clone the repository (if you haven't already)
git clone https://github.com/Mehulparekh144/goodforgo.git
cd tasks

# Build the binary
go build -o task-manager

# Make it executable
chmod +x task-manager

# Move to a directory in your PATH (optional)
sudo mv task-manager /usr/local/bin/
```

#### Windows

```powershell
# Clone the repository (if you haven't already)
git clone https://github.com/Mehulparekh144/goodforgo.git
cd tasks

# Build the binary
go build -o task-manager.exe
```

## File Storage

Tasks are stored in a local file under `./files/tasks.txt`

## Contributing

Feel free to submit issues and pull requests!

## License

MIT License
