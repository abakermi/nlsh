# Natural Language Shell (nlsh)

[![Build Status](https://github.com/abakermi/nlsh/actions/workflows/release.yml/badge.svg)](https://github.com/abakermi/nlsh/actions/workflows/release.yml)
[![Release](https://img.shields.io/github/v/release/abakermi/nlsh)](https://github.com/abakermi/nlsh/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/abakermi/nlsh)](https://github.com/abakermi/nlsh)

<img src="./resources//play.gif" width="500" />
A command-line tool that converts natural language instructions into shell commands using OpenAI's GPT model. Simply describe what you want to do in plain English, and nlsh will generate and execute the appropriate shell command.

## Features

- 🧠 Natural language to shell command conversion
- 🛡️ Built-in safety checks for dangerous commands
- ⚙️ Configurable settings via `.nlshrc`
- 🎨 Colored output for better readability
- 📝 Command history and context awareness
- 🔄 Interactive and single command modes
- 🔒 Confirmation for potentially dangerous operations

## Prerequisites

- Go 1.21 or later
- OpenAI API key

## Installation

### Option 1: Quick Install

Install directly using curl:
```bash
curl -fsSL https://raw.githubusercontent.com/abakermi/nlsh/master/install.sh | bash
```

### Option 2: Go Install

```bash
go install github.com/abakermi/nlsh@latest
```
### Option 3: Manual Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/abakermi/nlsh.git
   cd nlsh
   ```
2. Set your OpenAI API key as an environment variable:
   ```bash
   export OPENAI_API_KEY='your-api-key-here'
   ```
3. Run the installation script:
   ```bash
   ./install.sh
   ```
4. Restart your terminal or source your shell configuration:
   ```bash
   source ~/.zshrc  # or source ~/.bashrc
   ```

## Usage

### Set your OpenAI API key:

```bash
export OPENAI_API_KEY='your-api-key-here'
```
### Interactive Mode

```bash
nlsh
```

### Single Command Mode
```bash
nlsh "list all files in current directory"
```
## Examples
```bash
# List files
nlsh "show me all hidden files"

# Git operations
nlsh "commit all changes with message 'update readme'"

# Docker operations
nlsh "show all running containers"
```

## Safety Features
- Command confirmation before execution
- Configurable allowed/denied commands
- Pattern-based command filtering
- Protection against dangerous operations

## License

This project is open source and available under the MIT License.