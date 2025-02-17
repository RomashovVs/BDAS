package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

// Тесты для obfuscateString
func TestObfuscateString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Ifmmp"},
		{"123", "234"},
		{"ABC", "BCD"},
		{"", ""},
		{"Hello, World!", "Ifmmp- Xpsme\""},
	}

	for _, test := range tests {
		result := obfuscateString(test.input)
		if result != test.expected {
			t.Errorf("obfuscateString(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Тесты для deobfuscateString
func TestDeobfuscateString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Ifmmp", "Hello"},
		{"234", "123"},
		{"BCD", "ABC"},
		{"", ""},
		{"Ifmmp-!Xpsme\"", "Hello, World!"},
	}

	for _, test := range tests {
		result := deobfuscateString(test.input)
		if result != test.expected {
			t.Errorf("deobfuscateString(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Тесты для obfuscateText
func TestObfuscateText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Ifmmp"},
		{"<name>John</name>", "<name>Kpio</name>"},
		{"<email>john@example.com</email>", "<email>kpioAfybnqmf/dpn</email>"},
		{"", ""},
	}

	for _, test := range tests {
		result := obfuscateText(test.input)
		if result != test.expected {
			t.Errorf("obfuscateText(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Тесты для deobfuscateText
func TestDeobfuscateText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Ifmmp", "Hello"},
		{"<name>Kpio</name>", "<name>John</name>"},
		{"<email>kpioAfybnqmf/dpn</email>", "<email>john@example.com</email>"},
		{"", ""},
	}

	for _, test := range tests {
		result := deobfuscateText(test.input)
		if result != test.expected {
			t.Errorf("deobfuscateText(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Интеграционные тесты для obfuscateXML и deobfuscateXML
func TestObfuscateAndDeobfuscateXML(t *testing.T) {
	input := `
<document>
  <user id="1">
    <name>John Doe</name>
    <email>john.doe@example.com</email>
    <address>
      <street>123 Main St</street>
      <city>New York</city>
      <zip>10001</zip>
    </address>
  </user>
</document>
`

	// Обфускация
	obfuscated, err := obfuscateXML([]byte(input))
	if err != nil {
		t.Fatalf("obfuscateXML failed: %v", err)
	}

	// Деобфускация
	deobfuscated, err := deobfuscateXML(obfuscated)
	if err != nil {
		t.Fatalf("deobfuscateXML failed: %v", err)
	}

	// Сравнение с оригиналом
	var original, restored Document
	if err := xml.Unmarshal([]byte(input), &original); err != nil {
		t.Fatalf("xml.Unmarshal failed: %v", err)
	}
	if err := xml.Unmarshal(deobfuscated, &restored); err != nil {
		t.Fatalf("xml.Unmarshal failed: %v", err)
	}

	// Проверка, что структуры совпадают
	if !xmlEqual(original, restored) {
		t.Errorf("deobfuscated XML does not match original")
	}
}

// Вспомогательная функция для сравнения двух XML-документов
func xmlEqual(a, b Document) bool {
	aXML, _ := xml.Marshal(a)
	bXML, _ := xml.Marshal(b)
	return strings.TrimSpace(string(aXML)) == strings.TrimSpace(string(bXML))
}
