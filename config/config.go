package config

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

type env struct {
	// Default data files
	QuestionsFile string
	DeltasFile    string
	InfoLogFile   string
	ErrorLogFile  string
	SettingsFile  string
	// Default page size for pagination
	PageSize int
	MaxDelta int
	// Due Priority List: The default top-K due and upcoming questions to show in summary
	TopKDue      int
	TopKUpcoming int
	// Due Priority List: Scoring formula weights
	ImportanceWeight    float64
	OverdueWeight       float64
	FamiliarityWeight   float64
	ReviewPenaltyWeight float64
	EasePenaltyWeight   float64
	// SRS settings
	OverduePenalty    bool
	OverdueLimit      int
	RandomizeInterval bool
}

var (
	envInstance *env
	once        sync.Once
	defaultEnv  = &env{
		// Default data files
		QuestionsFile: "questions.default.json",
		DeltasFile:    "deltas.default.json",
		InfoLogFile:   "info.default.log",
		ErrorLogFile:  "error.default.log",
		SettingsFile:  "settings.default.json",
		// Default page size for pagination
		PageSize: 5,  // Default page size for pagination
		MaxDelta: 50, // Maximum number of deltas to keep
		// Due Priority List: The default top-K due and upcoming questions to show in summary
		TopKDue:      10, // Top-K due questions to show in summary
		TopKUpcoming: 10, // Top-K upcoming questions to show in summary
		// Due Priority List: Scoring formula weights
		ImportanceWeight:    1.5,  // Prioritizes designated importance
		OverdueWeight:       0.5,  // Prioritizes items past their due date
		FamiliarityWeight:   3.0,  // Prioritizes historically difficult items
		ReviewPenaltyWeight: -1.5, // De-prioritizes questions seen many times (prevents leeching)
		EasePenaltyWeight:   -1.0, // De-prioritizes "easier" questions to focus on struggles
		// SRS settings
		OverduePenalty:    false, // Enable/disable overdue penalty
		OverdueLimit:      7,     // Days after which overdue questions are at risk of penalty
		RandomizeInterval: true,  // Enable/disable randomized interval
	}
)

func Env() *env {
	once.Do(func() {
		envInstance = defaultEnv

		// Load environment variable overrides first
		envInstance.loadFromEnvironment()

		// Then load user settings file (which can override env vars)
		if err := envInstance.loadJSONFromFile(envInstance, envInstance.SettingsFile); err != nil {
			panic(errors.Wrap(err, "Failed to load settings file"))
		}
	})
	return envInstance
}

func (e *env) Save() error {
	return e.saveJSONToFile(e, e.SettingsFile)
}

func (e *env) loadJSONFromFile(data interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return nil for non-existent file
		}
		return err
	}
	defer file.Close()

	// Check if file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return nil // Leave data in its zero state
	}

	if err := json.NewDecoder(file).Decode(data); err != nil {
		return err
	}
	return nil
}

func (e *env) saveJSONToFile(data interface{}, filename string) error {
	tempFile, err := os.CreateTemp("", "temp_*.json")
	if err != nil {
		return err
	}

	// Ensure cleanup on any error
	cleanup := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}
	defer cleanup()

	// Write JSON data
	enc := json.NewEncoder(tempFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	// Close temp file before rename
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Atomic replace - disable cleanup since we want to keep the file
	cleanup = func() {} // No-op
	if err := os.Rename(tempFile.Name(), filename); err != nil {
		return err
	}
	return nil
}

// loadFromEnvironment reads configuration from environment variables
func (e *env) loadFromEnvironment() {
	// Default data files
	if val := os.Getenv("LEETSOLV_QUESTIONS_FILE"); val != "" {
		e.QuestionsFile = val
	}
	if val := os.Getenv("LEETSOLV_DELTAS_FILE"); val != "" {
		e.DeltasFile = val
	}
	if val := os.Getenv("LEETSOLV_INFO_LOG_FILE"); val != "" {
		e.InfoLogFile = val
	}
	if val := os.Getenv("LEETSOLV_ERROR_LOG_FILE"); val != "" {
		e.ErrorLogFile = val
	}
	if val := os.Getenv("LEETSOLV_SETTINGS_FILE"); val != "" {
		e.SettingsFile = val
	}
	// SRS settings
	if val := os.Getenv("LEETSOLV_OVERDUE_PENALTY"); val != "" {
		if overduePenalty, err := strconv.ParseBool(val); err == nil {
			e.OverduePenalty = overduePenalty
		}
	}
	if val := os.Getenv("LEETSOLV_RANDOMIZE_INTERVAL"); val != "" {
		if randomizeInterval, err := strconv.ParseBool(val); err == nil {
			e.RandomizeInterval = randomizeInterval
		}
	}
}

// ResetToDefaults resets all configuration values to their defaults except for the data files
func (e *env) ResetToDefaults() {
	questionsFile := e.QuestionsFile
	deltasFile := e.DeltasFile
	infoLogFile := e.InfoLogFile
	errorLogFile := e.ErrorLogFile
	settingsFile := e.SettingsFile

	*e = *defaultEnv
	e.QuestionsFile = questionsFile
	e.DeltasFile = deltasFile
	e.InfoLogFile = infoLogFile
	e.ErrorLogFile = errorLogFile
	e.SettingsFile = settingsFile
}

// Validate checks if the current configuration is valid
func (e *env) Validate() error {
	if e.PageSize <= 0 {
		return errors.New("PageSize must be positive")
	}
	if e.MaxDelta <= 0 {
		return errors.New("MaxDelta must be positive")
	}
	if e.TopKDue <= 0 {
		return errors.New("TopKDue must be positive")
	}
	if e.TopKUpcoming <= 0 {
		return errors.New("TopKUpcoming must be positive")
	}
	if e.OverdueLimit <= 0 {
		return errors.New("OverdueLimit must be positive")
	}
	return nil
}
