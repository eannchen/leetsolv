package copy

import (
	"testing"
	"time"
)

// TestStruct is a simple struct for testing deep copy functionality
type TestStruct struct {
	ID       int
	Name     string
	Active   bool
	Created  time.Time
	Tags     []string
	Metadata map[string]interface{}
}

// TestNestedStruct tests nested struct copying
type TestNestedStruct struct {
	ID       int
	Details  TestStruct
	Children []TestNestedStruct
}

func TestDeepCopyGob_SimpleTypes(t *testing.T) {
	tests := []struct {
		name     string
		src      interface{}
		expected interface{}
	}{
		{
			name:     "int",
			src:      42,
			expected: 42,
		},
		{
			name:     "string",
			src:      "hello world",
			expected: "hello world",
		},
		{
			name:     "bool",
			src:      true,
			expected: true,
		},
		{
			name:     "float64",
			src:      3.14159,
			expected: 3.14159,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a destination of the same type as source
			switch tt.src.(type) {
			case int:
				var dst int
				err := DeepCopyGob(&dst, tt.src)
				if err != nil {
					t.Fatalf("DeepCopyGob failed: %v", err)
				}
				if dst != tt.expected {
					t.Errorf("DeepCopyGob() = %v, want %v", dst, tt.expected)
				}
			case string:
				var dst string
				err := DeepCopyGob(&dst, tt.src)
				if err != nil {
					t.Fatalf("DeepCopyGob failed: %v", err)
				}
				if dst != tt.expected {
					t.Errorf("DeepCopyGob() = %v, want %v", dst, tt.expected)
				}
			case bool:
				var dst bool
				err := DeepCopyGob(&dst, tt.src)
				if err != nil {
					t.Fatalf("DeepCopyGob failed: %v", err)
				}
				if dst != tt.expected {
					t.Errorf("DeepCopyGob() = %v, want %v", dst, tt.expected)
				}
			case float64:
				var dst float64
				err := DeepCopyGob(&dst, tt.src)
				if err != nil {
					t.Fatalf("DeepCopyGob failed: %v", err)
				}
				if dst != tt.expected {
					t.Errorf("DeepCopyGob() = %v, want %v", dst, tt.expected)
				}
			}
		})
	}
}

func TestDeepCopyGob_Slice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	var dst []int

	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Errorf("DeepCopyGob() slice length = %d, want %d", len(dst), len(src))
	}

	for i, v := range src {
		if dst[i] != v {
			t.Errorf("DeepCopyGob() slice[%d] = %d, want %d", i, dst[i], v)
		}
	}

	// Verify it's a deep copy by modifying the original
	src[0] = 999
	if dst[0] == 999 {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src affected dst")
	}
}

func TestDeepCopyGob_Map(t *testing.T) {
	src := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	var dst map[string]int

	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Errorf("DeepCopyGob() map length = %d, want %d", len(dst), len(src))
	}

	for k, v := range src {
		if dst[k] != v {
			t.Errorf("DeepCopyGob() map[%s] = %d, want %d", k, dst[k], v)
		}
	}

	// Verify it's a deep copy by modifying the original
	src["a"] = 999
	if dst["a"] == 999 {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src affected dst")
	}
}

func TestDeepCopyGob_Struct(t *testing.T) {
	src := TestStruct{
		ID:      1,
		Name:    "test",
		Active:  true,
		Created: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:    []string{"tag1", "tag2"},
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	var dst TestStruct
	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	// Verify all fields are copied correctly
	if dst.ID != src.ID {
		t.Errorf("DeepCopyGob() ID = %d, want %d", dst.ID, src.ID)
	}
	if dst.Name != src.Name {
		t.Errorf("DeepCopyGob() Name = %s, want %s", dst.Name, src.Name)
	}
	if dst.Active != src.Active {
		t.Errorf("DeepCopyGob() Active = %t, want %t", dst.Active, src.Active)
	}
	if !dst.Created.Equal(src.Created) {
		t.Errorf("DeepCopyGob() Created = %v, want %v", dst.Created, src.Created)
	}

	// Verify slice is deep copied
	if len(dst.Tags) != len(src.Tags) {
		t.Errorf("DeepCopyGob() Tags length = %d, want %d", len(dst.Tags), len(src.Tags))
	}
	for i, tag := range src.Tags {
		if dst.Tags[i] != tag {
			t.Errorf("DeepCopyGob() Tags[%d] = %s, want %s", i, dst.Tags[i], tag)
		}
	}

	// Verify map is deep copied
	if len(dst.Metadata) != len(src.Metadata) {
		t.Errorf("DeepCopyGob() Metadata length = %d, want %d", len(dst.Metadata), len(src.Metadata))
	}
	for k, v := range src.Metadata {
		if dst.Metadata[k] != v {
			t.Errorf("DeepCopyGob() Metadata[%s] = %v, want %v", k, dst.Metadata[k], v)
		}
	}

	// Verify it's a deep copy by modifying the original
	src.Tags[0] = "modified"
	src.Metadata["key1"] = "modified"
	if dst.Tags[0] == "modified" {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src.Tags affected dst.Tags")
	}
	if dst.Metadata["key1"] == "modified" {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src.Metadata affected dst.Metadata")
	}
}

func TestDeepCopyGob_NestedStruct(t *testing.T) {
	src := TestNestedStruct{
		ID: 1,
		Details: TestStruct{
			ID:      2,
			Name:    "nested",
			Active:  true,
			Created: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Children: []TestNestedStruct{
			{
				ID: 3,
				Details: TestStruct{
					ID:      4,
					Name:    "child1",
					Active:  false,
					Created: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	var dst TestNestedStruct
	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	// Verify nested struct is copied correctly
	if dst.Details.ID != src.Details.ID {
		t.Errorf("DeepCopyGob() nested Details.ID = %d, want %d", dst.Details.ID, src.Details.ID)
	}
	if dst.Details.Name != src.Details.Name {
		t.Errorf("DeepCopyGob() nested Details.Name = %s, want %s", dst.Details.Name, src.Details.Name)
	}

	// Verify nested slice is copied correctly
	if len(dst.Children) != len(src.Children) {
		t.Errorf("DeepCopyGob() nested Children length = %d, want %d", len(dst.Children), len(src.Children))
	}
	if dst.Children[0].Details.Name != src.Children[0].Details.Name {
		t.Errorf("DeepCopyGob() nested Children[0].Details.Name = %s, want %s",
			dst.Children[0].Details.Name, src.Children[0].Details.Name)
	}

	// Verify it's a deep copy by modifying the original
	src.Details.Name = "modified"
	src.Children[0].Details.Name = "modified"
	if dst.Details.Name == "modified" {
		t.Error("DeepCopyGob() did not create a deep copy, modifying nested src affected dst")
	}
	if dst.Children[0].Details.Name == "modified" {
		t.Error("DeepCopyGob() did not create a deep copy, modifying nested src affected dst")
	}
}

func TestDeepCopyGob_EmptyValues(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		src := []int{}
		var dst []int
		err := DeepCopyGob(&dst, src)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		if len(dst) != 0 {
			t.Errorf("DeepCopyGob() empty slice length = %d, want 0", len(dst))
		}
	})

	t.Run("empty map", func(t *testing.T) {
		src := map[string]int{}
		var dst map[string]int
		err := DeepCopyGob(&dst, src)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		if len(dst) != 0 {
			t.Errorf("DeepCopyGob() empty map length = %d, want 0", len(dst))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var src []int
		var dst []int
		err := DeepCopyGob(&dst, src)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		// Gob converts nil slices to empty slices, which is expected behavior
		if len(dst) != 0 {
			t.Errorf("DeepCopyGob() nil slice length = %d, want 0", len(dst))
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var src map[string]int
		var dst map[string]int
		err := DeepCopyGob(&dst, src)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		// Gob converts nil maps to empty maps, which is expected behavior
		if len(dst) != 0 {
			t.Errorf("DeepCopyGob() nil map length = %d, want 0", len(dst))
		}
	})

	t.Run("empty struct", func(t *testing.T) {
		src := TestStruct{}
		var dst TestStruct
		err := DeepCopyGob(&dst, src)
		if err != nil {
			t.Fatalf("DeepCopyGob failed: %v", err)
		}
		// Verify empty struct is copied correctly
		if dst.ID != 0 || dst.Name != "" || dst.Active != false {
			t.Errorf("DeepCopyGob() empty struct = %+v, want zero values", dst)
		}
	})
}

func TestDeepCopyGob_ErrorCases(t *testing.T) {
	// Test with nil source
	var dst string
	err := DeepCopyGob(&dst, nil)
	if err == nil {
		t.Error("DeepCopyGob() should return error when src is nil")
	}

	// Test with unregistered type that can't be encoded
	unencodable := make(chan int)
	err = DeepCopyGob(&dst, unencodable)
	if err == nil {
		t.Error("DeepCopyGob() should return error when src can't be encoded")
	}
}

func TestDeepCopyGob_PointerTypes(t *testing.T) {
	value := 42
	src := &value
	var dst *int

	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	if dst == nil {
		t.Fatal("DeepCopyGob() dst should not be nil")
	}

	if *dst != *src {
		t.Errorf("DeepCopyGob() *dst = %d, want %d", *dst, *src)
	}

	// Verify it's a deep copy by modifying the original
	*src = 999
	if *dst == 999 {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src affected dst")
	}
}

func TestDeepCopyGob_InterfaceTypes(t *testing.T) {
	src := []interface{}{"string", 42, true, 3.14}
	var dst []interface{}

	err := DeepCopyGob(&dst, src)
	if err != nil {
		t.Fatalf("DeepCopyGob failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Errorf("DeepCopyGob() interface slice length = %d, want %d", len(dst), len(src))
	}

	for i, v := range src {
		if dst[i] != v {
			t.Errorf("DeepCopyGob() interface slice[%d] = %v, want %v", i, dst[i], v)
		}
	}

	// Verify it's a deep copy by modifying the original
	src[0] = "modified"
	if dst[0] == "modified" {
		t.Error("DeepCopyGob() did not create a deep copy, modifying src affected dst")
	}
}
