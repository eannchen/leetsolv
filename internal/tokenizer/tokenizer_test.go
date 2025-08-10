package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenize_EmptyString(t *testing.T) {
	result := Tokenize("")
	expected := []string{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WhitespaceOnly(t *testing.T) {
	result := Tokenize("   \t\n\r  ")
	expected := []string{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"   \\t\\n\\r  \") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_SingleWord(t *testing.T) {
	result := Tokenize("hello")
	expected := []string{"hello"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_MultipleWords(t *testing.T) {
	result := Tokenize("hello world")
	expected := []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello world\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_MixedCase(t *testing.T) {
	result := Tokenize("Hello World")
	expected := []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"Hello World\") returned %v, expected %v", result, expected)
	}

	result = Tokenize("HELLO WORLD")
	expected = []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"HELLO WORLD\") returned %v, expected %v", result, expected)
	}

	result = Tokenize("HeLLo WoRlD")
	expected = []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"HeLLo WoRlD\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithNumbers(t *testing.T) {
	result := Tokenize("hello123 world456")
	expected := []string{"hello123", "world456"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello123 world456\") returned %v, expected %v", result, expected)
	}

	result = Tokenize("123hello 456world")
	expected = []string{"123hello", "456world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"123hello 456world\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithPunctuation(t *testing.T) {
	result := Tokenize("hello, world!")
	expected := []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello, world!\") returned %v, expected %v", result, expected)
	}

	result = Tokenize("hello.world;test:case")
	expected = []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello.world;test:case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithSpecialCharacters(t *testing.T) {
	result := Tokenize("hello@world#test$case")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello@world#test$case\") returned %v, expected %v", result, expected)
	}

	result = Tokenize("hello&world|test+case")
	expected = []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello&world|test+case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithUnderscores(t *testing.T) {
	result := Tokenize("hello_world test_case")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello_world test_case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithHyphens(t *testing.T) {
	result := Tokenize("hello-world test-case")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello-world test-case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithParentheses(t *testing.T) {
	result := Tokenize("hello(world) test[case]")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello(world) test[case]\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithQuotes(t *testing.T) {
	result := Tokenize("hello\"world\" 'test'case")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello\\\"world\\\" 'test'case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithMultipleSpaces(t *testing.T) {
	result := Tokenize("hello    world")
	expected := []string{"hello", "world"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello    world\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithTabsAndNewlines(t *testing.T) {
	result := Tokenize("hello\tworld\ntest\r\ncase")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello\\tworld\\ntest\\r\\ncase\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithMixedDelimiters(t *testing.T) {
	result := Tokenize("hello,world;test:case")
	expected := []string{"hello", "world", "test", "case"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello,world;test:case\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithUnicode(t *testing.T) {
	result := Tokenize("café naïve résumé")
	expected := []string{"café", "naïve", "résumé"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"café naïve résumé\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithEmojis(t *testing.T) {
	result := Tokenize("hello 😀 world 🌍 test")
	expected := []string{"hello", "world", "test"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"hello 😀 world 🌍 test\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithURLs(t *testing.T) {
	result := Tokenize("visit https://example.com for more info")
	expected := []string{"visit", "https", "example", "com", "for", "more", "info"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"visit https://example.com for more info\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithEmail(t *testing.T) {
	result := Tokenize("contact me at user@example.com")
	expected := []string{"contact", "me", "at", "user", "example", "com"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"contact me at user@example.com\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_WithLeetCodeProblem(t *testing.T) {
	result := Tokenize("Two Sum - LeetCode Problem #1")
	expected := []string{"two", "sum", "leetcode", "problem", "1"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"Two Sum - LeetCode Problem #1\") returned %v, expected %v", result, expected)
	}
}

func TestTokenize_ComplexText(t *testing.T) {
	text := "This is a complex test case with multiple delimiters: commas, semicolons; periods. And numbers like 123 and 456!"
	result := Tokenize(text)
	expected := []string{"this", "is", "a", "complex", "test", "case", "with", "multiple", "delimiters", "commas", "semicolons", "periods", "and", "numbers", "like", "123", "and", "456"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Tokenize(\"%s\") returned %v, expected %v", text, result, expected)
	}
}

func TestTokenize_Consistency(t *testing.T) {
	// Test that the same input always produces the same output
	input := "Hello World Test Case"
	result1 := Tokenize(input)
	result2 := Tokenize(input)

	if !reflect.DeepEqual(result1, result2) {
		t.Errorf("Tokenize is not consistent: %v vs %v", result1, result2)
	}
}

func TestTokenize_NoModificationOfInput(t *testing.T) {
	// Test that the input string is not modified
	input := "Hello World"
	original := input
	_ = Tokenize(input)

	if input != original {
		t.Error("Tokenize modified the input string")
	}
}
