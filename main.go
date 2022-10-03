package main

import (
	"bufio"
	"combiner/comb"
	"errors"
	"fmt"
	"os"
)

const defaultTragetFileName string = "target.txt"

func main() {

	config := comb.CombinerConfig{}

	// Verify args and prepare config struct
	args := os.Args[1:]
	verifyArgs(args, &config)

	lines := make(map[string]string)

	// Read the first supplied file and add its lines to a map
	readFileAndAppend(config.FirstFilePath, lines)

	if config.HasStdInInput {
		scanner := bufio.NewScanner(os.Stdin)
		scanAndAppend(scanner, lines)
	} else {
		readFileAndAppend(config.SecondFilePath, lines)
	}

	dFileWriter, err := os.Create(config.TargetFilePath)
	if err != nil {
		fmt.Println("Error occured while creating taregt file: ", err)
		os.Exit(1)
	}

	for _, line := range lines {
		dFileWriter.WriteString(line + "\n")
	}

}

func scanAndAppend(scanner *bufio.Scanner, lines map[string]string) {
	for scanner.Scan() {
		lineText := scanner.Text()
		if _, ok := lines[lineText]; ok {
			fmt.Println("Duplicate line found, ignoring it: ", lineText)
		} else {
			lines[lineText] = lineText
		}
	}
}

func verifyArgs(args []string, config *comb.CombinerConfig) {

	// check if there is stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		config.HasStdInInput = true
	} else {
		config.HasStdInInput = false
	}

	if !config.HasStdInInput {
		if len(args) < 2 {
			fmt.Println("Please supply two files locations")
			os.Exit(1)
		}

		if !fileExists(args[0]) || !fileExists(args[1]) {
			fmt.Println("One of the supplied files is not found")
			os.Exit(1)
		}

		config.FirstFilePath = args[0]
		config.SecondFilePath = args[1]

	}

	if config.HasStdInInput {
		if len(args) < 1 {
			fmt.Println("Please supply a file path")
			os.Exit(1)
		}

		if !fileExists(args[0]) {
			fmt.Println("One of the supplied files is not found")
			os.Exit(1)
		}
		config.FirstFilePath = args[0]

	}

	config.TargetFilePath = defaultTragetFileName
}

func readFileAndAppend(path string, lines map[string]string) {
	fileReader, err := os.Open(path)
	if err != nil {
		fmt.Println("Error coccured while reading file: ", err)
	}
	fileScanner := bufio.NewScanner(fileReader)
	fileScanner.Split(bufio.ScanLines)
	scanAndAppend(fileScanner, lines)
	fileReader.Close()
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Couldn't find file: " + path)
		return false
	}
	return true
}
