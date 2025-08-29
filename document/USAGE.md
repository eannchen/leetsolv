# Usage

## Interactive Mode
```bash
# Start interactive session
leetsolv

# You'll see the prompt:
leetsolv ❯
```

## Command Line Mode
```bash
# List all questions
leetsolv list

# Search for problems with filters
leetsolv search tree --familiarity=3 --importance=2 --due-only

# Get problem details
leetsolv detail 123

# Check status
leetsolv status

# Add new problem
leetsolv add https://leetcode.com/problems/example

# After re-solving it, update to schedule the next review
leetsolv upsert https://leetcode.com/problems/example
```

## Available Commands

| Command   | Aliases               | Description                                     |
| --------- | --------------------- | ----------------------------------------------- |
| `list`    | `ls`                  | List all questions with pagination              |
| `search`  | `s`                   | Search questions by keywords (supports filters) |
| `detail`  | `get`                 | Get detailed information about a question       |
| `status`  | `stat`                | Show summary of due and upcoming questions      |
| `upsert`  | `add`                 | Add or update a question                        |
| `remove`  | `rm`, `delete`, `del` | Delete a question                               |
| `undo`    | `back`                | Undo the last action                            |
| `history` | `hist`, `log`         | Show action history                             |
| `setting` | `config`, `cfg`       | View and modify application settings            |
| `version` | `ver`, `v`            | Show application version information            |
| `help`    | `h`                   | Show help information                           |
| `clear`   | `cls`                 | Clear the screen                                |
| `quit`    | `q`, `exit`           | Exit the application                            |


## Search Command Filters

The `search` command lets you search by keywords (in **URL** or **note**) and refine results using filters.

**Syntax:**
```bash
search [keywords...] [filters...]
```

**Filters:**

| Filter             | Description                       |
| ------------------ | --------------------------------- |
| `--familiarity=N`  | Filter by familiarity level (1-5) |
| `--importance=N`   | Filter by importance level (1-4)  |
| `--review-count=N` | Filter by review count            |
| `--due-only`       | Only show due questions           |