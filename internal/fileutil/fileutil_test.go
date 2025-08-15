package fileutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Test data structures
type TestData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ComplexTestData struct {
	Items []TestData `json:"items"`
	Meta  struct {
		Version string `json:"version"`
		Count   int    `json:"count"`
	} `json:"meta"`
}

func TestNewJSONFileUtil(t *testing.T) {
	util := NewJSONFileUtil()
	if util == nil {
		t.Error("NewJSONFileUtil should return a non-nil instance")
	}
}

func TestJSONFileUtil_Load_NonExistentFile(t *testing.T) {
	util := NewJSONFileUtil()
	var data TestData
	
	// Try to load from a non-existent file
	err := util.Load(&data, "/path/that/does/not/exist.json")
	
	// Should not return an error for non-existent files
	if err != nil {
		t.Errorf("Load should not return error for non-existent file, got: %v", err)
	}
	
	// Data should remain in zero state
	if data.Name != "" || data.Value != 0 {
		t.Errorf("Data should remain in zero state, got: %+v", data)
	}
}

func TestJSONFileUtil_Load_EmptyFile(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create an empty temporary file
	tempFile, err := os.CreateTemp("", "test_empty_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	var data TestData
	err = util.Load(&data, tempFile.Name())
	
	// Should not return an error for empty files
	if err != nil {
		t.Errorf("Load should not return error for empty file, got: %v", err)
	}
	
	// Data should remain in zero state
	if data.Name != "" || data.Value != 0 {
		t.Errorf("Data should remain in zero state, got: %+v", data)
	}
}

func TestJSONFileUtil_Load_ValidFile(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file with valid JSON
	tempFile, err := os.CreateTemp("", "test_valid_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write test data to file
	testData := TestData{Name: "test", Value: 42}
	encoder := json.NewEncoder(tempFile)
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()
	
	// Load the data
	var loadedData TestData
	err = util.Load(&loadedData, tempFile.Name())
	
	if err != nil {
		t.Errorf("Load should not return error for valid file, got: %v", err)
	}
	
	// Verify loaded data matches original
	if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
		t.Errorf("Loaded data doesn't match original. Expected: %+v, Got: %+v", testData, loadedData)
	}
}

func TestJSONFileUtil_Load_ComplexData(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file with complex JSON
	tempFile, err := os.CreateTemp("", "test_complex_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write complex test data to file
	testData := ComplexTestData{
		Items: []TestData{
			{Name: "item1", Value: 10},
			{Name: "item2", Value: 20},
		},
		Meta: struct {
			Version string `json:"version"`
			Count   int    `json:"count"`
		}{
			Version: "1.0.0",
			Count:   2,
		},
	}
	
	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()
	
	// Load the data
	var loadedData ComplexTestData
	err = util.Load(&loadedData, tempFile.Name())
	
	if err != nil {
		t.Errorf("Load should not return error for valid complex file, got: %v", err)
	}
	
	// Verify loaded data matches original
	if len(loadedData.Items) != len(testData.Items) {
		t.Errorf("Items length mismatch. Expected: %d, Got: %d", len(testData.Items), len(loadedData.Items))
	}
	
	if loadedData.Meta.Version != testData.Meta.Version {
		t.Errorf("Meta version mismatch. Expected: %s, Got: %s", testData.Meta.Version, loadedData.Meta.Version)
	}
	
	if loadedData.Meta.Count != testData.Meta.Count {
		t.Errorf("Meta count mismatch. Expected: %d, Got: %d", testData.Meta.Count, loadedData.Meta.Count)
	}
}

func TestJSONFileUtil_Load_InvalidJSON(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file with invalid JSON
	tempFile, err := os.CreateTemp("", "test_invalid_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write invalid JSON
	if _, err := tempFile.WriteString(`{"name": "test", "value": invalid}`); err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}
	tempFile.Close()
	
	var data TestData
	err = util.Load(&data, tempFile.Name())
	
	// Should return an error for invalid JSON
	if err == nil {
		t.Error("Load should return error for invalid JSON")
	}
}

func TestJSONFileUtil_Save_ValidData(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file path
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_save_valid.json")
	defer os.Remove(tempFile)
	
	// Save test data
	testData := TestData{Name: "saved_test", Value: 123}
	err := util.Save(testData, tempFile)
	
	if err != nil {
		t.Errorf("Save should not return error for valid data, got: %v", err)
	}
	
	// Verify file was created and contains correct data
	var loadedData TestData
	err = util.Load(&loadedData, tempFile)
	if err != nil {
		t.Errorf("Failed to load saved data: %v", err)
	}
	
	if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
		t.Errorf("Saved data doesn't match original. Expected: %+v, Got: %+v", testData, loadedData)
	}
}

func TestJSONFileUtil_Save_ComplexData(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file path
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_save_complex.json")
	defer os.Remove(tempFile)
	
	// Save complex test data
	testData := ComplexTestData{
		Items: []TestData{
			{Name: "saved_item1", Value: 100},
			{Name: "saved_item2", Value: 200},
			{Name: "saved_item3", Value: 300},
		},
		Meta: struct {
			Version string `json:"version"`
			Count   int    `json:"count"`
		}{
			Version: "2.0.0",
			Count:   3,
		},
	}
	
	err := util.Save(testData, tempFile)
	if err != nil {
		t.Errorf("Save should not return error for valid complex data, got: %v", err)
	}
	
	// Verify file was created and contains correct data
	var loadedData ComplexTestData
	err = util.Load(&loadedData, tempFile)
	if err != nil {
		t.Errorf("Failed to load saved complex data: %v", err)
	}
	
	// Verify all fields match
	if len(loadedData.Items) != len(testData.Items) {
		t.Errorf("Items length mismatch. Expected: %d, Got: %d", len(testData.Items), len(loadedData.Items))
	}
	
	for i, item := range testData.Items {
		if loadedData.Items[i].Name != item.Name || loadedData.Items[i].Value != item.Value {
			t.Errorf("Item %d mismatch. Expected: %+v, Got: %+v", i, item, loadedData.Items[i])
		}
	}
}

func TestJSONFileUtil_Save_AtomicOperation(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file path
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_atomic.json")
	defer os.Remove(tempFile)
	
	// First, save some initial data
	initialData := TestData{Name: "initial", Value: 1}
	err := util.Save(initialData, tempFile)
	if err != nil {
		t.Fatalf("Failed to save initial data: %v", err)
	}
	
	// Verify initial data was saved
	var loadedData TestData
	err = util.Load(&loadedData, tempFile)
	if err != nil {
		t.Fatalf("Failed to load initial data: %v", err)
	}
	
	if loadedData.Name != initialData.Name || loadedData.Value != initialData.Value {
		t.Errorf("Initial data mismatch. Expected: %+v, Got: %+v", initialData, loadedData)
	}
	
	// Now save new data - this should atomically replace the file
	newData := TestData{Name: "updated", Value: 2}
	err = util.Save(newData, tempFile)
	if err != nil {
		t.Errorf("Failed to save updated data: %v", err)
	}
	
	// Verify new data replaced old data
	err = util.Load(&loadedData, tempFile)
	if err != nil {
		t.Fatalf("Failed to load updated data: %v", err)
	}
	
	if loadedData.Name != newData.Name || loadedData.Value != newData.Value {
		t.Errorf("Updated data mismatch. Expected: %+v, Got: %+v", newData, loadedData)
	}
}

func TestJSONFileUtil_Save_InvalidPath(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Try to save to an invalid path (directory that doesn't exist)
	invalidPath := "/path/that/does/not/exist/file.json"
	testData := TestData{Name: "test", Value: 42}
	
	err := util.Save(testData, invalidPath)
	
	// Should return an error for invalid path
	if err == nil {
		t.Error("Save should return error for invalid path")
	}
}

func TestJSONFileUtil_SaveAndLoad_RoundTrip(t *testing.T) {
	util := NewJSONFileUtil()
	
	// Create a temporary file path
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_roundtrip.json")
	defer os.Remove(tempFile)
	
	// Test multiple round trips with different data
	testCases := []TestData{
		{Name: "first", Value: 1},
		{Name: "second", Value: 2},
		{Name: "third", Value: 3},
		{Name: "", Value: 0}, // Test zero values
		{Name: "special chars: !@#$%^&*()", Value: -999},
	}
	
	for i, testData := range testCases {
		// Save data
		err := util.Save(testData, tempFile)
		if err != nil {
			t.Errorf("Test case %d: Save failed: %v", i, err)
			continue
		}
		
		// Load data
		var loadedData TestData
		err = util.Load(&loadedData, tempFile)
		if err != nil {
			t.Errorf("Test case %d: Load failed: %v", i, err)
			continue
		}
		
		// Verify data matches
		if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
			t.Errorf("Test case %d: Data mismatch. Expected: %+v, Got: %+v", i, testData, loadedData)
		}
	}
}

func TestJSONFileUtil_Interface_Compliance(t *testing.T) {
	// Verify that JSONFileUtil implements FileUtil interface
	var _ FileUtil = (*JSONFileUtil)(nil)
	
	// Test that we can use it through the interface
	var util FileUtil = NewJSONFileUtil()
	
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_interface.json")
	defer os.Remove(tempFile)
	
	testData := TestData{Name: "interface_test", Value: 999}
	
	// Save through interface
	err := util.Save(testData, tempFile)
	if err != nil {
		t.Errorf("Save through interface failed: %v", err)
	}
	
	// Load through interface
	var loadedData TestData
	err = util.Load(&loadedData, tempFile)
	if err != nil {
		t.Errorf("Load through interface failed: %v", err)
	}
	
	// Verify data
	if loadedData.Name != testData.Name || loadedData.Value != testData.Value {
		t.Errorf("Interface data mismatch. Expected: %+v, Got: %+v", testData, loadedData)
	}
}
