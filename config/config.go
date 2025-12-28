// Package config implements the configuration for the leetsolv application.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eannchen/leetsolv/internal/copy"
	"github.com/eannchen/leetsolv/internal/fileutil"
)

var (
	defaultConfig *Config

	// Environment variable loaders
	envLoaders = []struct {
		Key   string
		Apply func(e *Config, val string)
	}{
		{"LEETSOLV_QUESTIONS_FILE", func(e *Config, v string) { e.QuestionsFile = v }},
		{"LEETSOLV_DELTAS_FILE", func(e *Config, v string) { e.DeltasFile = v }},
		{"LEETSOLV_INFO_LOG_FILE", func(e *Config, v string) { e.InfoLogFile = v }},
		{"LEETSOLV_ERROR_LOG_FILE", func(e *Config, v string) { e.ErrorLogFile = v }},
		{"LEETSOLV_SETTINGS_FILE", func(e *Config, v string) { e.SettingsFile = v }},
		{"LEETSOLV_RANDOMIZE_INTERVAL", func(e *Config, v string) {
			if b, err := strconv.ParseBool(v); err == nil {
				e.RandomizeInterval = b
			}
		}},
		{"LEETSOLV_OVERDUE_PENALTY", func(e *Config, v string) {
			if b, err := strconv.ParseBool(v); err == nil {
				e.OverduePenalty = b
			}
		}},
		{"LEETSOLV_OVERDUE_LIMIT", func(e *Config, v string) {
			if i, err := strconv.Atoi(v); err == nil {
				e.OverdueLimit = i
			}
		}},
	}

	// Settings registry (for configurable settings)
	settingsRegistry = map[string]SettingDefinition{
		"randomizeinterval": {
			Name:        "RandomizeInterval",
			Type:        "bool",
			Description: "Enable/disable randomized interval",
			Validator: func(valueStr string) (any, error) {
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					return boolValue, nil
				}
				return nil, errors.New("RandomizeInterval must be a boolean value")
			},
			Getter: func(e *Config) any {
				return e.RandomizeInterval
			},
			Setter: func(e *Config, value any) error {
				if boolValue, ok := value.(bool); ok {
					e.RandomizeInterval = boolValue
					return nil
				}
				return errors.New("RandomizeInterval must be a boolean value")
			},
		},
		"overduepenalty": {
			Name:        "OverduePenalty",
			Type:        "bool",
			Description: "Enable/disable overdue penalty",
			Validator: func(valueStr string) (any, error) {
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					return boolValue, nil
				}
				return nil, errors.New("OverduePenalty must be a boolean value")
			},
			Getter: func(e *Config) any {
				return e.OverduePenalty
			},
			Setter: func(e *Config, value any) error {
				if boolValue, ok := value.(bool); ok {
					e.OverduePenalty = boolValue
					return nil
				}
				return errors.New("OverduePenalty must be a boolean value")
			},
		},
		"overduelimit": {
			Name:        "OverdueLimit",
			Type:        "int",
			Unit:        "days",
			Description: "Days after which overdue questions are at risk of penalty",
			Validator: func(valueStr string) (any, error) {
				if intValue, err := strconv.Atoi(valueStr); err == nil {
					return intValue, nil
				}
				return nil, errors.New("OverdueLimit must be an integer value")
			},
			Getter: func(e *Config) any {
				return e.OverdueLimit
			},
			Setter: func(e *Config, value any) error {
				if intValue, ok := value.(int); ok {
					e.OverdueLimit = intValue
					return e.validate()
				}
				return errors.New("OverdueLimit must be an integer value")
			},
		},
	}
)

// initDefaultConfig initializes the default configuration with proper file paths
func initDefaultConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".leetsolv")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	defaultConfig = &Config{
		// Default data files with absolute paths
		QuestionsFile: filepath.Join(configDir, "questions.json"),
		DeltasFile:    filepath.Join(configDir, "deltas.json"),
		InfoLogFile:   filepath.Join(configDir, "info.log"),
		ErrorLogFile:  filepath.Join(configDir, "error.log"),
		SettingsFile:  filepath.Join(configDir, "settings.json"),
		// Pagination settings
		Paginator: Paginator{
			PageSize: 5,
		},
		// Delta settings
		Delta: Delta{
			MaxDelta: 25,
		},
		// Due Priority List settings
		DuePriority: DuePriority{
			TopKDue:             10,   // Top-K due questions to show in summary
			TopKUpcoming:        10,   // Top-K upcoming questions to show in summary
			ImportanceWeight:    1.5,  // Prioritizes designated importance
			OverdueWeight:       0.5,  // Prioritizes items past their due date
			FamiliarityWeight:   3.0,  // Prioritizes historically difficult items
			ReviewPenaltyWeight: -1.5, // De-prioritizes questions seen many times (prevents leeching)
			EasePenaltyWeight:   -1.0, // De-prioritizes "easier" questions to focus on struggles
		},
		// SRS settings
		SRS: SRS{
			RandomizeInterval: true,  // Enable/disable randomized interval
			OverduePenalty:    false, // Enable/disable overdue penalty
			OverdueLimit:      7,     // Days after which overdue questions are at risk of penalty
		},
	}

	return nil
}

// SettingDefinition defines a configurable setting
type SettingDefinition struct {
	Name        string
	Type        string // "bool", "int", "float64", "string"
	Unit        string
	Description string
	Validator   func(string) (any, error)
	Getter      func(*Config) any
	Setter      func(*Config, any) error
}

type Paginator struct {
	// Default page size for pagination
	PageSize int `json:"pageSize"`
}

type Delta struct {
	// Maximum number of deltas to keep
	MaxDelta int `json:"maxDelta"`
}

type DuePriority struct {
	// The default top-K due and upcoming questions to show in summary
	TopKDue      int `json:"topKDue"`
	TopKUpcoming int `json:"topKUpcoming"`
	// Scoring formula weights
	ImportanceWeight    float64 `json:"importanceWeight"`
	OverdueWeight       float64 `json:"overdueWeight"`
	FamiliarityWeight   float64 `json:"familiarityWeight"`
	ReviewPenaltyWeight float64 `json:"reviewPenaltyWeight"`
	EasePenaltyWeight   float64 `json:"easePenaltyWeight"`
}

type SRS struct {
	RandomizeInterval bool `json:"randomizeInterval"`
	OverduePenalty    bool `json:"overduePenalty"`
	OverdueLimit      int  `json:"overdueLimit"`
}

type Config struct {
	// Dependency injection
	file fileutil.FileUtil
	// Default data files
	QuestionsFile string `json:"questionsFile"`
	DeltasFile    string `json:"deltasFile"`
	InfoLogFile   string `json:"infoLogFile"`
	ErrorLogFile  string `json:"errorLogFile"`
	SettingsFile  string `json:"settingsFile"`
	// Pagination settings
	Paginator
	// Delta settings
	Delta
	// Due Priority List settings
	DuePriority
	// SRS settings
	SRS
}

func NewConfig(file fileutil.FileUtil) (*Config, error) {
	// Initialize default config with proper paths
	if err := initDefaultConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize default configuration: %v", err)
	}

	config := &Config{file: file}

	// Load default configuration first
	if err := config.loadFromDefault(); err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %v", err)
	}

	// Load environment variable overrides first
	if err := config.loadFromEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %v", err)
	}

	// Then load user settings file (which can override env vars)
	if err := config.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load settings file: %w", err)
	}

	return config, nil
}

func (e *Config) Save() error {
	return e.file.Save(e, e.SettingsFile)
}

func (e *Config) loadFromDefault() error {
	if err := copy.DeepCopyGob(e, defaultConfig); err != nil {
		return err
	}
	return e.validate()
}

func (e *Config) loadFromFile() error {
	if err := e.file.Load(e, e.SettingsFile); err != nil {
		return err
	}
	return e.validate()
}

// loadFromEnvironment reads configuration from environment variables
func (e *Config) loadFromEnvironment() error {
	for _, loader := range envLoaders {
		if val := os.Getenv(loader.Key); val != "" {
			loader.Apply(e, val)
		}
	}
	return e.validate()
}

// Validate checks if the current configuration is valid
func (e *Config) validate() error {
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
func (e *Config) GetSettingsRegistry() map[string]SettingDefinition {
	return settingsRegistry
}

// GetSettingValue retrieves a configurable setting value by name
func (e *Config) GetSettingValue(settingName string) (any, error) {
	setting, exists := settingsRegistry[strings.ToLower(settingName)]
	if !exists {
		return nil, fmt.Errorf("unknown setting: %s", settingName)
	}
	return setting.Getter(e), nil
}

// SetSettingValue sets a configurable setting value by name
func (e *Config) SetSettingValue(settingName string, value any) error {
	setting, exists := settingsRegistry[strings.ToLower(settingName)]
	if !exists {
		return fmt.Errorf("unknown setting: %s", settingName)
	}
	return setting.Setter(e, value)
}

// GetSettingInfo returns information about a configurable setting
func (e *Config) GetSettingInfo(settingName string) (*SettingDefinition, error) {
	setting, exists := settingsRegistry[strings.ToLower(settingName)]
	if !exists {
		return nil, fmt.Errorf("unknown setting: %s", settingName)
	}
	return &setting, nil
}
