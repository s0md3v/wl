package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func isLower(char string) bool {
	// checks if char is lowercase
	return unicode.IsLower([]rune(char)[0])
}

func isUpper(char string) bool {
	// checks if char is uppercase
	return unicode.IsUpper([]rune(char)[0])
}

func isMixed(str string) bool {
	// checks if string is mixed case
	lower, upper := false, false
	for _, char := range str {
		if isLower(string(char)) {
			lower = true
		} else if isUpper(string(char)) {
			upper = true
		}
		if lower && upper {
			return true
		}
	}
	return false
}

func add(str string, del string, i int) string {
	// concates delimeter if i (index) is non-zero
	if i == 0 {
		return str
	}
	return del + str
}

func detectFooBar(fooBar string) (string, string) {
	// detects the casing style of given string
	del := ""
	casing := ""
	if isMixed(fooBar) {
		if isLower(string(fooBar[0])) {
			casing = "c"
		} else {
			casing = "p"
		}
	} else if isLower(string(fooBar[0])) {
		casing = "l"
	} else {
		casing = "u"
	}
	if strings.Contains(fooBar, "-") {
		del = "-"
	} else if strings.Contains(fooBar, "_") {
		del = "_"
	} else if strings.Contains(fooBar, ".") {
		del = "."
	}
	return del, casing
}

func transform(parts []string, del string, casing string) string {
	// combines list of strings to form a string with given casing style
	str := ""
	if len(parts) == 1 {
		if casing == "l" {
			return strings.ToLower(parts[0])
		} else if casing == "u" {
			return strings.ToUpper(parts[0])
		}
		return parts[0]
	}
	for i, part := range parts {
		if casing == "l" {
			str += add(strings.ToLower(part), del, i)
		} else if casing == "u" {
			str += add(strings.ToUpper(part), del, i)
		} else if casing == "c" {
			if i == 0 {
				str += add(strings.ToLower(part), del, i)
			} else {
				str += add(strings.Title(strings.ToLower(part)), del, i)
			}
		} else if casing == "p" {
			str += add(strings.Title(strings.ToLower(part)), del, i)
		}
	}
	return str
}

func handle(str string) []string {
	// breaks down a string into array of 'words'
	if strings.Contains(str, "-") {
		return strings.Split(str, "-")
	} else if strings.Contains(str, "_") {
		return strings.Split(str, "_")
	} else if strings.Contains(str, ".") {
		return strings.Split(str, ".")
	}
	parts := []string{}
	temp := ""
	if isMixed(str) {
		for _, char := range str {
			if !isUpper(string(char)) {
				temp += string(char)
			} else {
				if temp != "" {
					parts = append(parts, temp)
				}
				temp = string(char)
			}
		}
		parts = append(parts, temp)
		return parts
	}
	return []string{str}
}

func start(casing string, inputStream io.Reader, outputStream io.Writer) {
	del, casing := detectFooBar(casing)
	scanner := bufio.NewScanner(inputStream)
	writer := bufio.NewWriter(outputStream)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		parts := handle(line)
		transformedLine := transform(parts, del, casing) + "\n"
		writer.WriteString(transformedLine)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func main() {
	casing := flag.String("c", "", "casing style (required)")
	inputFile := flag.String("i", "", "input file (default stdin)")
	outputFile := flag.String("o", "", "output file (default stdout)")
	flag.Parse()

	if *casing == "" {
		flag.Usage()
		return
	}

	outputStream := os.Stdout
	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		outputStream = file
	}

	inputStream := os.Stdin
	if *inputFile != "" {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		inputStream = file
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read standard input: %v\n", err)
			os.Exit(1)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			flag.Usage()
			return
		}
	}

	start(*casing, inputStream, outputStream)
}
