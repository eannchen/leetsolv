# Data Migration Script

This script migrates the old `questions.json` format to the new `QuestionStore` format used by the Go application.

## Usage

1. Make sure you have the old `questions.json` file in the project root directory
2. Run the migration script:

```bash
go run cmd/migrate/main.go
```

## What the script does

1. **Reads the old data**: Parses the existing `questions.json` file
2. **Converts data types**:
   - Converts familiarity from integers (0-4) to enum values (VeryHard-VeryEasy)
   - Converts importance from integers (1-3) to enum values (LowImportance-HighImportance)
   - Parses timestamp strings into proper `time.Time` objects
3. **Creates new structure**: Builds the new `QuestionStore` with:
   - `Questions` map indexed by ID
   - `URLIndex` map for quick URL lookups
   - `MaxID` for tracking the highest question ID
4. **Saves migrated data**: Outputs to `questions_new.json` in the new format

## Data Mapping

### Familiarity
- Old: 0 → New: VeryHard
- Old: 1 → New: Hard
- Old: 2 → New: Medium
- Old: 3 → New: Easy
- Old: 4 → New: VeryEasy

### Importance
- Old: 1 → New: MediumImportance (preserves value 1)
- Old: 2 → New: HighImportance (preserves value 2)
- Old: 3 → New: CriticalImportance (preserves value 3)

## Output

The script will:
- Create `questions_new.json` with the migrated data
- Display migration statistics
- Show familiarity and importance distributions

## After Migration

1. Review the generated `questions_new.json` file
2. If satisfied, replace your existing questions file:
   ```bash
   mv questions_new.json questions.json
   ```
3. Update your application configuration to use the new file format

## Error Handling

The script includes error handling for:
- Missing or invalid timestamp formats (uses current time as fallback)
- Invalid familiarity/importance values (uses Medium as default)
- File I/O errors

All warnings are logged but don't stop the migration process.