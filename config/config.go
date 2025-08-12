package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"leetsolv/internal/copy"
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
	RandomizeInterval bool
	OverduePenalty    bool
	OverdueLimit      int
}

// SettingDefinition defines a configurable setting
type SettingDefinition struct {
	Name         string
	Type         string // "bool", "int", "float64", "string"
	Description  string
	DefaultValue interface{}
	Validator    func(interface{}) error
	Getter       func(*env) interface{}
	Setter       func(*env, interface{}) error
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
		RandomizeInterval: true,  // Enable/disable randomized interval
		OverduePenalty:    false, // Enable/disable overdue penalty
		OverdueLimit:      7,     // Days after which overdue questions are at risk of penalty
	}
)

func Env() *env {
	once.Do(func() {
		envInstance = &env{}
		if err := copy.DeepCopyGob(envInstance, defaultEnv); err != nil {
			panic(errors.Wrap(err, "Failed to copy default environment"))
		}

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
func (e *env) ResetToDefaults() error {
	questionsFile := e.QuestionsFile
	deltasFile := e.DeltasFile
	infoLogFile := e.InfoLogFile
	errorLogFile := e.ErrorLogFile
	settingsFile := e.SettingsFile

	if err := copy.DeepCopyGob(e, defaultEnv); err != nil {
		return errors.Wrap(err, "Failed to copy default environment")
	}
	e.QuestionsFile = questionsFile
	e.DeltasFile = deltasFile
	e.InfoLogFile = infoLogFile
	e.ErrorLogFile = errorLogFile
	e.SettingsFile = settingsFile
	return nil
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

// GetSettingsRegistry returns the registry of all configurable settings
func GetSettingsRegistry() map[string]SettingDefinition {
	return map[string]SettingDefinition{
		"randomizeinterval": {
			Name:         "RandomizeInterval",
			Type:         "bool",
			Description:  "Enable/disable randomized interval",
			DefaultValue: true,
			Validator: func(value interface{}) error {
				if _, ok := value.(bool); !ok {
					return errors.New("RandomizeInterval must be a boolean value")
				}
				return nil
			},
			Getter: func(e *env) interface{} {
				return e.RandomizeInterval
			},
			Setter: func(e *env, value interface{}) error {
				if boolValue, ok := value.(bool); ok {
					e.RandomizeInterval = boolValue
					return nil
				}
				return errors.New("RandomizeInterval must be a boolean value")
			},
		},
		"overduepenalty": {
			Name:         "OverduePenalty",
			Type:         "bool",
			Description:  "Enable/disable overdue penalty",
			DefaultValue: false,
			Validator: func(value interface{}) error {
				if _, ok := value.(bool); !ok {
					return errors.New("OverduePenalty must be a boolean value")
				}
				return nil
			},
			Getter: func(e *env) interface{} {
				return e.OverduePenalty
			},
			Setter: func(e *env, value interface{}) error {
				if boolValue, ok := value.(bool); ok {
					e.OverduePenalty = boolValue
					return nil
				}
				return errors.New("OverduePenalty must be a boolean value")
			},
		},
		"overduelimit": {
			Name:         "OverdueLimit",
			Type:         "int",
			Description:  "Days after which overdue questions are at risk of penalty",
			DefaultValue: 7,
			Validator: func(value interface{}) error {
				if intValue, ok := value.(int); ok {
					if intValue <= 0 {
						return errors.New("OverdueLimit must be a positive integer")
					}
					return nil
				}
				return errors.New("OverdueLimit must be an integer value")
			},
			Getter: func(e *env) interface{} {
				return e.OverdueLimit
			},
			Setter: func(e *env, value interface{}) error {
				if intValue, ok := value.(int); ok {
					if intValue <= 0 {
						return errors.New("OverdueLimit must be a positive integer")
					}
					e.OverdueLimit = intValue
					return nil
				}
				return errors.New("OverdueLimit must be an integer value")
			},
		},
	}
}

// GetSettingValue retrieves a setting value by name
func (e *env) GetSettingValue(settingName string) (interface{}, error) {
	registry := GetSettingsRegistry()
	setting, exists := registry[strings.ToLower(settingName)]
	if !exists {
		return nil, fmt.Errorf("Unknown setting: %s", settingName)
	}
	return setting.Getter(e), nil
}

// SetSettingValue sets a setting value by name
func (e *env) SetSettingValue(settingName string, value interface{}) error {
	registry := GetSettingsRegistry()
	setting, exists := registry[strings.ToLower(settingName)]
	if !exists {
		return fmt.Errorf("Unknown setting: %s", settingName)
	}

	if err := setting.Validator(value); err != nil {
		return err
	}

	return setting.Setter(e, value)
}

// GetSettingInfo returns information about a setting
func GetSettingInfo(settingName string) (*SettingDefinition, error) {
	registry := GetSettingsRegistry()
	setting, exists := registry[strings.ToLower(settingName)]
	if !exists {
		return nil, fmt.Errorf("Unknown setting: %s", settingName)
	}
	return &setting, nil
}
