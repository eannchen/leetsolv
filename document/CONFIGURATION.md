# Configuration

You can configure **LeetSolv** in two ways:

1. **Environment variables** – convenient for temporary or deployment-level overrides.
2. **JSON settings file** (`$HOME/.leetsolv/settings.json`) – persistent configuration you can edit manually.

Both methods map to the same internal configuration.
- Environment variables follow `UPPERCASE_SNAKE_CASE` naming.
- JSON fields follow `camelCase` naming.

For example:
- Env var: `LEETSOLV_RANDOMIZE_INTERVAL=true`
- JSON: `"randomizeInterval": true`

If both are provided, the **JSON settings file takes priority** over environment variables.

## File Paths

| Env Variable              | JSON field      | Default                          | Description         |
| ------------------------- | --------------- | -------------------------------- | ------------------- |
| `LEETSOLV_QUESTIONS_FILE` | `questionsFile` | `$HOME/.leetsolv/questions.json` | Questions data file |
| `LEETSOLV_DELTAS_FILE`    | `deltasFile`    | `$HOME/.leetsolv/deltas.json`    | Change history file |
| `LEETSOLV_INFO_LOG_FILE`  | `infoLogFile`   | `$HOME/.leetsolv/info.log`       | Info log file       |
| `LEETSOLV_ERROR_LOG_FILE` | `errorLogFile`  | `$HOME/.leetsolv/error.log`      | Error log file      |
| `LEETSOLV_SETTINGS_FILE`  | `settingsFile`  | `$HOME/.leetsolv/settings.json`  | Config JSON file    |


## SM-2 Algorithm Settings

| Env Variable                  | JSON field          | Default | Description                                    |
| ----------------------------- | ------------------- | ------- | ---------------------------------------------- |
| `LEETSOLV_RANDOMIZE_INTERVAL` | `randomizeInterval` | `true`  | Enable/disable interval randomization          |
| `LEETSOLV_OVERDUE_PENALTY`    | `overduePenalty`    | `false` | Enable/disable overdue penalty system          |
| `LEETSOLV_OVERDUE_LIMIT`      | `overdueLimit`      | `7`     | Days after which overdue questions get penalty |


## Due Priority Scoring Settings

| Env Variable                     | JSON field            | Default | Description                    |
| -------------------------------- | --------------------- | ------- | ------------------------------ |
| `LEETSOLV_TOP_K_DUE`             | `topKDue`             | `10`    | Top due questions to show      |
| `LEETSOLV_TOP_K_UPCOMING`        | `topKUpcoming`        | `10`    | Top upcoming questions to show |
| `LEETSOLV_IMPORTANCE_WEIGHT`     | `importanceWeight`    | `1.5`   | Weight for problem importance  |
| `LEETSOLV_OVERDUE_WEIGHT`        | `overdueWeight`       | `0.5`   | Weight for overdue problems    |
| `LEETSOLV_FAMILIARITY_WEIGHT`    | `familiarityWeight`   | `3.0`   | Weight for familiarity level   |
| `LEETSOLV_REVIEW_PENALTY_WEIGHT` | `reviewPenaltyWeight` | `-1.5`  | Penalty for high review count  |
| `LEETSOLV_EASE_PENALTY_WEIGHT`   | `easePenaltyWeight`   | `-1.0`  | Penalty for easy problems      |


## Other Settings

| Env Variable         | JSON field | Default | Description             |
| -------------------- | ---------- | ------- | ----------------------- |
| `LEETSOLV_PAGE_SIZE` | `pageSize` | `5`     | Questions per page      |
| `LEETSOLV_MAX_DELTA` | `maxDelta` | `50`    | Maximum history entries |

## Example: Environment Variables

```bash
export LEETSOLV_RANDOMIZE_INTERVAL=false
export LEETSOLV_PAGE_SIZE=20
```

## Example: JSON Settings File

```json
{
    "randomizeInterval": false,
    "pageSize": 20
}
```
