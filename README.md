# LeetSolv

A sophisticated spaced repetition system for LeetCode problem management, built in Go with an intelligent scheduling algorithm and comprehensive CLI interface.

## 🎯 Overview

LeetSolv is a command-line tool designed to help developers systematically review and master LeetCode problems using spaced repetition principles. It implements the SM-2 algorithm with custom adaptations for coding interview preparation, featuring intelligent prioritization based on problem difficulty, importance, and review history.

## ✨ Features

### Core Functionality
- **Spaced Repetition**: SM-2 algorithm implementation with custom intervals based on problem importance
- **Intelligent Scheduling**: Dynamic review scheduling considering familiarity, importance, and memory retention
- **Priority Scoring**: Advanced scoring system that prioritizes questions based on multiple factors
- **Search & Filtering**: Powerful search capabilities with multiple filter options
- **History Tracking**: Complete audit trail of all changes with undo functionality

### Problem Management
- **Add/Update Problems**: Easy problem entry with URL and notes
- **Importance Levels**: 4-tier importance system (Low, Medium, High, Critical)
- **Familiarity Tracking**: 5-level familiarity scale (VeryHard → VeryEasy)
- **Memory Assessment**: Track how well you remember each problem
- **Due Date Management**: Automatic calculation of next review dates

### CLI Interface
- **Interactive Mode**: Full-featured interactive CLI with command history
- **Batch Mode**: Execute commands directly from command line arguments
- **Alias Support**: Multiple command aliases for convenience
- **Pagination**: Efficient handling of large problem sets
- **Clear Output**: Well-formatted, readable command output

## 🏗️ Architecture

### Project Structure
```
leetsolv/
├── core/           # Core domain models and business logic
├── usecase/        # Application use cases and orchestration
├── handler/        # Request handling and user interaction
├── command/        # CLI command implementations
├── storage/        # Data persistence layer
├── internal/       # Internal utilities and helpers
├── config/         # Configuration management
└── main.go         # Application entry point
```

### Key Components

#### Core Domain (`core/`)
- **Question Model**: Central data structure with all problem metadata
- **SM2Scheduler**: Spaced repetition algorithm implementation
- **Action Tracking**: Delta-based change history system

#### Use Cases (`usecase/`)
- **QuestionUseCase**: Main business logic for problem management
- **Search & Filtering**: Advanced query capabilities
- **Priority Calculation**: Intelligent scoring algorithms

#### Storage (`storage/`)
- **File-based Storage**: JSON-based data persistence
- **Delta Tracking**: Change history with rollback support
- **Atomic Operations**: Safe file operations with error handling

#### Command System (`command/`)
- **Command Registry**: Extensible command system
- **Handler Integration**: Clean separation of concerns
- **Alias Support**: Multiple command names for convenience

## 🚀 Installation

### Prerequisites
- Go 1.24.4 or later
- Git

### Build from Source
```bash
# Clone the repository
git clone <repository-url>
cd leetsolv

# Build the application
go build -o leetsolv .

# Make executable
chmod +x leetsolv

# Run (optional)
./leetsolv
```

### Quick Start
```bash
# Development mode (uses test files)
make dev

# Production mode (uses production files)
make prod

# Show available make targets
make help
```

## 📖 Usage

### Interactive Mode
```bash
# Start interactive session
./leetsolv

# You'll see the prompt:
leetsolv ❯
```

### Command Line Mode
```bash
# List all questions
./leetsolv list

# Search for problems
./leetsolv search "binary tree"

# Get problem details
./leetsolv get 123

# Check status
./leetsolv status

# Add new problem
./leetsolv add "https://leetcode.com/problems/example"
```

### Available Commands

| Command   | Aliases               | Description                                |
| --------- | --------------------- | ------------------------------------------ |
| `list`    | `ls`                  | List all questions with pagination         |
| `search`  | `s`                   | Search questions by keywords               |
| `get`     | `detail`              | Get detailed information about a question  |
| `status`  | `stat`                | Show summary of due and upcoming questions |
| `add`     | `upsert`              | Add or update a question                   |
| `remove`  | `rm`, `delete`, `del` | Delete a question                          |
| `undo`    | `back`                | Undo the last action                       |
| `history` | `hist`, `log`         | Show action history                        |
| `help`    | `h`                   | Show help information                      |
| `clear`   | `cls`                 | Clear the screen                           |
| `quit`    | `q`, `exit`           | Exit the application                       |

## ⚙️ Configuration

### Environment Variables

| Variable                  | Default               | Description         |
| ------------------------- | --------------------- | ------------------- |
| `LEETSOLV_QUESTIONS_FILE` | `questions.test.json` | Questions data file |
| `LEETSOLV_DELTAS_FILE`    | `deltas.test.json`    | Change history file |
| `LEETSOLV_INFO_LOG_FILE`  | `info.test.log`       | Info log file       |
| `LEETSOLV_ERROR_LOG_FILE` | `error.test.log`      | Error log file      |

### Scoring Weights

| Variable                         | Default | Description                   |
| -------------------------------- | ------- | ----------------------------- |
| `LEETSOLV_IMPORTANCE_WEIGHT`     | `1.5`   | Weight for problem importance |
| `LEETSOLV_OVERDUE_WEIGHT`        | `0.5`   | Weight for overdue problems   |
| `LEETSOLV_FAMILIARITY_WEIGHT`    | `3.0`   | Weight for difficulty level   |
| `LEETSOLV_REVIEW_PENALTY_WEIGHT` | `-1.5`  | Penalty for high review count |
| `LEETSOLV_EASE_PENALTY_WEIGHT`   | `-1.0`  | Penalty for easy problems     |

### Other Settings

| Variable                  | Default | Description                    |
| ------------------------- | ------- | ------------------------------ |
| `LEETSOLV_PAGE_SIZE`      | `5`     | Questions per page             |
| `LEETSOLV_MAX_DELTA`      | `50`    | Maximum history entries        |
| `LEETSOLV_TOP_K_DUE`      | `10`    | Top due questions to show      |
| `LEETSOLV_TOP_K_UPCOMING` | `10`    | Top upcoming questions to show |

## 🔧 Development

### Project Setup
```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run specific test
go test ./usecase -v

# Run integration tests
go test ./usecase -tags=integration
```

### Architecture Principles

#### Clean Architecture
- **Separation of Concerns**: Clear boundaries between layers
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each component has one clear purpose

#### Design Patterns
- **Command Pattern**: Extensible command system
- **Strategy Pattern**: Pluggable scheduling algorithms
- **Repository Pattern**: Abstracted data access
- **Factory Pattern**: Dependency creation and management

#### Error Handling
- **Structured Errors**: Coded error types with context
- **Graceful Degradation**: Application continues working despite errors
- **Comprehensive Logging**: Detailed error tracking and debugging

### Testing Strategy
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end workflow testing
- **Test Utilities**: Reusable test helpers and mocks
- **Coverage Goals**: High test coverage for critical paths

### Code Quality
- **Go Best Practices**: Following Go idioms and conventions
- **Error Handling**: Proper error propagation and logging
- **Documentation**: Comprehensive code documentation
- **Consistent Formatting**: Go fmt and linting compliance

## 📊 Data Model

### Question Structure
```go
type Question struct {
    ID           int         // Unique identifier
    URL          string      // LeetCode problem URL
    Note         string      // Personal notes
    Familiarity  Familiarity // Difficulty level (VeryHard → VeryEasy)
    Importance   Importance  // Priority level (Low → Critical)
    LastReviewed time.Time   // Last review timestamp
    NextReview   time.Time   // Next scheduled review
    ReviewCount  int         // Number of reviews completed
    EaseFactor   float64     // SM-2 ease factor
    UpdatedAt    time.Time   // Last modification time
    CreatedAt    time.Time   // Creation timestamp
}
```

### Scheduling Algorithm
The SM-2 scheduler adapts the standard spaced repetition algorithm:

1. **Base Intervals**: Different starting intervals based on importance
2. **Memory Multipliers**: Adjust intervals based on memory performance
3. **Familiarity Adjustments**: Early difficulty signals affect scheduling
4. **Ease Factor Management**: Dynamic adjustment of review intervals
5. **Maximum Limits**: Prevents excessively long intervals

## 🚀 Performance Features

### Efficient Data Structures
- **Priority Heaps**: Fast top-K queries for due/upcoming questions
- **Trie Search**: Efficient text search and filtering
- **Lazy Loading**: On-demand data loading and processing

### Memory Management
- **Streaming Operations**: Handle large datasets without memory issues
- **Pagination**: Efficient display of large question sets
- **Delta Compression**: Efficient storage of change history

## 🔒 Data Safety

### File Operations
- **Atomic Writes**: Safe file updates with temporary files
- **Backup Creation**: Automatic backup before major changes
- **Error Recovery**: Graceful handling of file corruption

### History Management
- **Complete Audit Trail**: Every change is recorded
- **Undo Capability**: Rollback any action
- **Delta Storage**: Efficient change tracking

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Code Standards
- Follow Go formatting standards (`go fmt`)
- Run linter checks (`golangci-lint`)
- Maintain test coverage
- Add documentation for new features

## 📝 License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## 🆘 Support

### Common Issues
- **File Permissions**: Ensure write access to data files
- **Data Corruption**: Use `undo` command or restore from backup
- **Performance**: Check file sizes and consider cleanup

### Getting Help
- Check existing issues in the repository
- Review the code documentation
- Run with verbose logging for debugging

---

**LeetSolv** - Master LeetCode problems with intelligent spaced repetition.