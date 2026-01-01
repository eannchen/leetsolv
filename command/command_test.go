package command

import (
	"bufio"
	"strings"
	"testing"
)

// MockHandler implements handler.Handler for testing
type MockHandler struct {
	listCalled    bool
	searchCalled  bool
	getCalled     bool
	statusCalled  bool
	upsertCalled  bool
	deleteCalled  bool
	undoCalled    bool
	helpCalled    bool
	clearCalled   bool
	quitCalled    bool
	historyCalled bool
	settingCalled bool
	versionCalled bool
	migrateCalled bool
	resetCalled   bool

	searchArgs  []string
	getArgs     string
	upsertArgs  string
	deleteArgs  string
	settingArgs []string
}

func (m *MockHandler) HandleList(scanner *bufio.Scanner) {
	m.listCalled = true
}

func (m *MockHandler) HandleSearch(scanner *bufio.Scanner, args []string) {
	m.searchCalled = true
	m.searchArgs = args
}

func (m *MockHandler) HandleGet(scanner *bufio.Scanner, target string) {
	m.getCalled = true
	m.getArgs = target
}

func (m *MockHandler) HandleStatus() {
	m.statusCalled = true
}

func (m *MockHandler) HandleUpsert(scanner *bufio.Scanner, rawURL string) {
	m.upsertCalled = true
	m.upsertArgs = rawURL
}

func (m *MockHandler) HandleDelete(scanner *bufio.Scanner, target string) {
	m.deleteCalled = true
	m.deleteArgs = target
}

func (m *MockHandler) HandleUndo(scanner *bufio.Scanner) {
	m.undoCalled = true
}

func (m *MockHandler) HandleHelp() {
	m.helpCalled = true
}

func (m *MockHandler) HandleClear() {
	m.clearCalled = true
}

func (m *MockHandler) HandleQuit() {
	m.quitCalled = true
}

func (m *MockHandler) HandleHistory() {
	m.historyCalled = true
}

func (m *MockHandler) HandleSetting(scanner *bufio.Scanner, args []string) {
	m.settingCalled = true
	m.settingArgs = args
}

func (m *MockHandler) HandleUnknown(command string) {
	// Not used in command tests
}

func (m *MockHandler) HandleVersion() {
	m.versionCalled = true
}

func (m *MockHandler) HandleMigrate(scanner *bufio.Scanner) {
	m.migrateCalled = true
}

func (m *MockHandler) HandleReset(scanner *bufio.Scanner) {
	m.resetCalled = true
}

func TestNewCommandRegistry(t *testing.T) {
	unknownHandler := func(command string) {
		// This handler is just for testing the constructor
	}

	registry := NewCommandRegistry(unknownHandler)

	if registry == nil {
		t.Fatal("NewCommandRegistry returned nil")
	}

	if registry.commands == nil {
		t.Error("commands map should not be nil")
	}

	if registry.unknownCommandHandler == nil {
		t.Error("unknownCommandHandler should not be nil")
	}
}

func TestCommandRegistry_Register(t *testing.T) {
	registry := NewCommandRegistry(func(command string) {})
	mockHandler := &MockHandler{}

	// Test case-sensitive registration
	registry.Register("Test", &ListCommand{Handler: mockHandler})

	// Should be stored in lowercase
	if _, exists := registry.commands["test"]; !exists {
		t.Error("Command should be registered in lowercase")
	}

	if _, exists := registry.commands["Test"]; exists {
		t.Error("Command should not be registered in original case")
	}
}

func TestCommandRegistry_Execute_ExistingCommand(t *testing.T) {
	registry := NewCommandRegistry(func(command string) {})
	mockHandler := &MockHandler{}

	// Register a command
	registry.Register("test", &ListCommand{Handler: mockHandler})

	// Execute the command
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := registry.Execute(scanner, "test", []string{})

	if quit {
		t.Error("ListCommand should not return quit=true")
	}

	if !mockHandler.listCalled {
		t.Error("Handler should have been called")
	}
}

func TestCommandRegistry_Execute_NonExistentCommand(t *testing.T) {
	unknownHandlerCalled := false
	unknownHandler := func(command string) {
		unknownHandlerCalled = true
	}

	registry := NewCommandRegistry(unknownHandler)

	// Execute non-existent command
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := registry.Execute(scanner, "nonexistent", []string{})

	if quit {
		t.Error("Non-existent command should not return quit=true")
	}

	if !unknownHandlerCalled {
		t.Error("Unknown command handler should have been called")
	}
}

func TestCommandRegistry_Execute_CaseInsensitive(t *testing.T) {
	registry := NewCommandRegistry(func(command string) {})
	mockHandler := &MockHandler{}

	// Register command in lowercase
	registry.Register("test", &ListCommand{Handler: mockHandler})

	// Execute with different cases
	scanner := bufio.NewScanner(strings.NewReader(""))

	// Test uppercase
	registry.Execute(scanner, "TEST", []string{})
	if !mockHandler.listCalled {
		t.Error("Command should be found regardless of case")
	}

	// Reset mock
	mockHandler.listCalled = false

	// Test mixed case
	registry.Execute(scanner, "TeSt", []string{})
	if !mockHandler.listCalled {
		t.Error("Command should be found regardless of case")
	}
}

func TestListCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &ListCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("ListCommand should not return quit=true")
	}

	if !mockHandler.listCalled {
		t.Error("Handler.HandleList should have been called")
	}
}

func TestSearchCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &SearchCommand{Handler: mockHandler}

	args := []string{"query1", "query2"}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("SearchCommand should not return quit=true")
	}

	if !mockHandler.searchCalled {
		t.Error("Handler.HandleSearch should have been called")
	}

	if len(mockHandler.searchArgs) != 2 {
		t.Errorf("Expected 2 search args, got %d", len(mockHandler.searchArgs))
	}

	if mockHandler.searchArgs[0] != "query1" || mockHandler.searchArgs[1] != "query2" {
		t.Error("Search args not passed correctly")
	}
}

func TestGetCommand_Execute_WithArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &GetCommand{Handler: mockHandler}

	args := []string{"123"}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("GetCommand should not return quit=true")
	}

	if !mockHandler.getCalled {
		t.Error("Handler.HandleGet should have been called")
	}

	if mockHandler.getArgs != "123" {
		t.Errorf("Expected target '123', got '%s'", mockHandler.getArgs)
	}
}

func TestGetCommand_Execute_WithoutArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &GetCommand{Handler: mockHandler}

	args := []string{}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("GetCommand should not return quit=true")
	}

	if !mockHandler.getCalled {
		t.Error("Handler.HandleGet should have been called")
	}

	if mockHandler.getArgs != "" {
		t.Errorf("Expected empty target, got '%s'", mockHandler.getArgs)
	}
}

func TestStatusCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &StatusCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("StatusCommand should not return quit=true")
	}

	if !mockHandler.statusCalled {
		t.Error("Handler.HandleStatus should have been called")
	}
}

func TestUpsertCommand_Execute_WithArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &UpsertCommand{Handler: mockHandler}

	args := []string{"https://leetcode.com/problems/test"}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("UpsertCommand should not return quit=true")
	}

	if !mockHandler.upsertCalled {
		t.Error("Handler.HandleUpsert should have been called")
	}

	if mockHandler.upsertArgs != "https://leetcode.com/problems/test" {
		t.Errorf("Expected URL 'https://leetcode.com/problems/test', got '%s'", mockHandler.upsertArgs)
	}
}

func TestUpsertCommand_Execute_WithoutArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &UpsertCommand{Handler: mockHandler}

	args := []string{}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("UpsertCommand should not return quit=true")
	}

	if !mockHandler.upsertCalled {
		t.Error("Handler.HandleUpsert should have been called")
	}

	if mockHandler.upsertArgs != "" {
		t.Errorf("Expected empty URL, got '%s'", mockHandler.upsertArgs)
	}
}

func TestDeleteCommand_Execute_WithArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &DeleteCommand{Handler: mockHandler}

	args := []string{"123"}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("DeleteCommand should not return quit=true")
	}

	if !mockHandler.deleteCalled {
		t.Error("Handler.HandleDelete should have been called")
	}

	if mockHandler.deleteArgs != "123" {
		t.Errorf("Expected target '123', got '%s'", mockHandler.deleteArgs)
	}
}

func TestDeleteCommand_Execute_WithoutArgs(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &DeleteCommand{Handler: mockHandler}

	args := []string{}
	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, args)

	if quit {
		t.Error("DeleteCommand should not return quit=true")
	}

	if !mockHandler.deleteCalled {
		t.Error("Handler.HandleDelete should have been called")
	}

	if mockHandler.deleteArgs != "" {
		t.Errorf("Expected empty target, got '%s'", mockHandler.deleteArgs)
	}
}

func TestUndoCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &UndoCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("UndoCommand should not return quit=true")
	}

	if !mockHandler.undoCalled {
		t.Error("Handler.HandleUndo should have been called")
	}
}

func TestHelpCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &HelpCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("HelpCommand should not return quit=true")
	}

	if !mockHandler.helpCalled {
		t.Error("Handler.HandleHelp should have been called")
	}
}

func TestClearCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &ClearCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("ClearCommand should not return quit=true")
	}

	if !mockHandler.clearCalled {
		t.Error("Handler.HandleClear should have been called")
	}
}

func TestQuitCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &QuitCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if !quit {
		t.Error("QuitCommand should return quit=true")
	}

	if !mockHandler.quitCalled {
		t.Error("Handler.HandleQuit should have been called")
	}
}

func TestHistoryCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &HistoryCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("HistoryCommand should not return quit=true")
	}

	if !mockHandler.historyCalled {
		t.Error("Handler.HandleHistory should have been called")
	}
}

func TestCommandRegistry_RegisterMultipleCommands(t *testing.T) {
	registry := NewCommandRegistry(func(command string) {})
	mockHandler := &MockHandler{}

	// Register multiple commands
	registry.Register("list", &ListCommand{Handler: mockHandler})
	registry.Register("search", &SearchCommand{Handler: mockHandler})
	registry.Register("get", &GetCommand{Handler: mockHandler})

	// Verify all commands are registered
	if len(registry.commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(registry.commands))
	}

	// Verify each command can be executed
	scanner := bufio.NewScanner(strings.NewReader(""))

	registry.Execute(scanner, "list", []string{})
	if !mockHandler.listCalled {
		t.Error("List command not executed")
	}

	mockHandler.listCalled = false
	registry.Execute(scanner, "search", []string{})
	if !mockHandler.searchCalled {
		t.Error("Search command not executed")
	}

	mockHandler.searchCalled = false
	registry.Execute(scanner, "get", []string{})
	if !mockHandler.getCalled {
		t.Error("Get command not executed")
	}
}

func TestCommandRegistry_ExecuteWithScanner(t *testing.T) {
	registry := NewCommandRegistry(func(command string) {})
	mockHandler := &MockHandler{}

	// Register a command that uses the scanner
	registry.Register("test", &ListCommand{Handler: mockHandler})

	// Create a scanner with some input
	input := "test input\nmore input"
	scanner := bufio.NewScanner(strings.NewReader(input))

	// Execute the command
	registry.Execute(scanner, "test", []string{})

	// Verify the command was executed
	if !mockHandler.listCalled {
		t.Error("Command should have been executed")
	}
}

func TestCommandInterfaceCompliance(t *testing.T) {
	// Test that all commands implement the Command interface
	var _ Command = &ListCommand{}
	var _ Command = &SearchCommand{}
	var _ Command = &GetCommand{}
	var _ Command = &StatusCommand{}
	var _ Command = &UpsertCommand{}
	var _ Command = &DeleteCommand{}
	var _ Command = &UndoCommand{}
	var _ Command = &HelpCommand{}
	var _ Command = &ClearCommand{}
	var _ Command = &QuitCommand{}
	var _ Command = &HistoryCommand{}
	var _ Command = &SettingCommand{}
	var _ Command = &VersionCommand{}
	var _ Command = &MigrateCommand{}
	var _ Command = &ResetCommand{}
}

func TestSettingCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &SettingCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{"randomizeinterval", "true"})

	if quit {
		t.Error("SettingCommand should not return quit=true")
	}

	if !mockHandler.settingCalled {
		t.Error("Handler.HandleSetting should have been called")
	}
}

func TestVersionCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &VersionCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("VersionCommand should not return quit=true")
	}

	if !mockHandler.versionCalled {
		t.Error("Handler.HandleVersion should have been called")
	}
}

func TestMigrateCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &MigrateCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("MigrateCommand should not return quit=true")
	}

	if !mockHandler.migrateCalled {
		t.Error("Handler.HandleMigrate should have been called")
	}
}

func TestResetCommand_Execute(t *testing.T) {
	mockHandler := &MockHandler{}
	command := &ResetCommand{Handler: mockHandler}

	scanner := bufio.NewScanner(strings.NewReader(""))
	quit := command.Execute(scanner, []string{})

	if quit {
		t.Error("ResetCommand should not return quit=true")
	}

	if !mockHandler.resetCalled {
		t.Error("Handler.HandleReset should have been called")
	}
}
