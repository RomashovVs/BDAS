package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Document struct {
	XMLName xml.Name `xml:"document"`
	Nodes   []Node   `xml:",any"`
}

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	InnerXML string     `xml:",innerxml"` // Захватываем всё содержимое узла
}

func obfuscateXML(data []byte) ([]byte, error) {
	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	obfuscateNode(&doc)
	return xml.MarshalIndent(doc, "", "  ")
}

func obfuscateNode(node *Document) {
	for i := range node.Nodes {
		n := &node.Nodes[i]
		// Обрабатываем атрибуты
		for j := range n.Attrs {
			n.Attrs[j].Value = obfuscateString(n.Attrs[j].Value)
		}
		// Обрабатываем текстовое содержимое
		n.InnerXML = obfuscateText(n.InnerXML)
	}
}

// obfuscateText обрабатывает текстовое содержимое узла
func obfuscateText(s string) string {
	// Разбираем InnerXML вручную, чтобы обработать только текстовые данные
	decoder := xml.NewDecoder(strings.NewReader(s))
	var result strings.Builder
	encoder := xml.NewEncoder(&result)

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch v := token.(type) {
		case xml.CharData:
			// Обрабатываем только текстовые данные
			encoder.EncodeToken(xml.CharData(obfuscateString(string(v))))
		case xml.StartElement, xml.EndElement, xml.Comment, xml.ProcInst, xml.Directive:
			// Записываем остальные токены без изменений
			encoder.EncodeToken(v)
		}
	}

	encoder.Flush()
	return result.String()
}

func obfuscateString(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r <= 126 && r != ' ' && r != '\n' && r != '\t' { // Игнорируем пробелы, \n и \t
			return r + 1
		}
		return r
	}, s)
}

func deobfuscateXML(data []byte) ([]byte, error) {
	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	deobfuscateNode(&doc)
	return xml.MarshalIndent(doc, "", "  ")
}

func deobfuscateNode(node *Document) {
	for i := range node.Nodes {
		n := &node.Nodes[i]
		// Обрабатываем атрибуты
		for j := range n.Attrs {
			n.Attrs[j].Value = deobfuscateString(n.Attrs[j].Value)
		}
		// Обрабатываем текстовое содержимое
		n.InnerXML = deobfuscateText(n.InnerXML)
	}
}

// deobfuscateText обрабатывает текстовое содержимое узла
func deobfuscateText(s string) string {
	// Разбираем InnerXML вручную, чтобы обработать только текстовые данные
	decoder := xml.NewDecoder(strings.NewReader(s))
	var result strings.Builder
	encoder := xml.NewEncoder(&result)

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch v := token.(type) {
		case xml.CharData:
			// Обрабатываем только текстовые данные
			encoder.EncodeToken(xml.CharData(deobfuscateString(string(v))))
		case xml.StartElement, xml.EndElement, xml.Comment, xml.ProcInst, xml.Directive:
			// Записываем остальные токены без изменений
			encoder.EncodeToken(v)
		}
	}

	encoder.Flush()
	return result.String()
}

func deobfuscateString(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 33 && r <= 127 && r != ' ' && r != '\n' && r != '\t' { // Игнорируем пробелы, \n и \t
			return r - 1
		}
		return r
	}, s)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: go run main.go <input.xml> <output.xml> [--deobfuscate]")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	deobfuscate := len(os.Args) > 3 && os.Args[3] == "--deobfuscate"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	var result []byte
	if deobfuscate {
		result, err = deobfuscateXML(data)
	} else {
		result, err = obfuscateXML(data)
	}

	if err != nil {
		fmt.Println("Ошибка обработки XML:", err)
		return
	}

	if err := os.WriteFile(outputFile, result, 0644); err != nil {
		fmt.Println("Ошибка записи файла:", err)
		return
	}

	fmt.Println("Операция завершена успешно.")
}
