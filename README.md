# LeetSolv

A spaced repetition system for LeetCode problems with intelligent priority scheduling.

## Priority Scoring System

LeetSolv uses a configurable priority scoring algorithm to determine which questions should be reviewed first. The system automatically prioritizes questions based on multiple factors to ensure efficient learning.

### Scoring Formula

The priority score is calculated using the following formula:

```
Priority Score = (ImportanceWeight × Importance) + (OverdueWeight × Overdue Days) + (FamiliarityWeight × Difficulty) + (ReviewPenaltyWeight × Review Count) + (EasePenaltyWeight × Ease Factor)
```

Where:
- **Importance**: 1-4 scale (Low=1, Medium=2, High=3, Critical=4)
- **Overdue Days**: Days past due date (minimum 0)
- **Difficulty**: 4 - Familiarity (VeryEasy=0, VeryHard=4)
- **Review Count**: Number of times reviewed
- **Ease Factor**: SM-2 algorithm ease factor

### Default Weights

The default scoring weights are:
- `ImportanceWeight`: 1.5 (prioritizes designated importance)
- `OverdueWeight`: 0.5 (prioritizes items past their due date)
- `FamiliarityWeight`: 3.0 (prioritizes historically difficult items)
- `ReviewPenaltyWeight`: -1.5 (de-prioritizes questions seen many times)
- `EasePenaltyWeight`: -1.0 (de-prioritizes "easier" questions)

### Customizing the Scoring Formula

You can customize the scoring weights by setting environment variables:

```bash
# Override individual weights
export LEETSOLV_IMPORTANCE_WEIGHT=2.0
export LEETSOLV_OVERDUE_WEIGHT=1.0
export LEETSOLV_FAMILIARITY_WEIGHT=4.0
export LEETSOLV_REVIEW_PENALTY_WEIGHT=-2.0
export LEETSOLV_EASE_PENALTY_WEIGHT=-1.5

# Run LeetSolv with custom weights
./leetsolv
```

### Viewing Current Scoring

Use the `status` command to see the current scoring formula and weights:

```bash
./leetsolv status
```

This will display the current formula with the actual weights being used, making it easy to verify your configuration.

## Usage

[Add usage instructions here]

## Configuration

[Add configuration instructions here]