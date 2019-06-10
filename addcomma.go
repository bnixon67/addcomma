package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// if error, then panic
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// process each file
func processFile(fileName string) {
	fmt.Printf("%s\n", fileName)

	// input file
	inFile, err := os.Open(fileName)
	checkError(err)
	defer inFile.Close()

	// output file
	outFile, err := os.Create(filepath.Base(fileName))
	checkError(err)
	defer outFile.Close()

	// loop thru lines of input file via Scanner
	inScanner := bufio.NewScanner(inFile)
	for inScanner.Scan() {
		// add comma to end of line, using Windows CRLF
		outFile.WriteString(inScanner.Text())
		outFile.WriteString(",\r\n")
	}

	// check for any non-EOF errors during Scan
	if err := inScanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getDirs(baseDirName string) []string {

	// open base directory
	dir, err := os.Open(baseDirName)
	checkError(err)
	defer dir.Close()

	// read all files in directory
	files, err := dir.Readdir(0)
	checkError(err)

	var dirs []string

	// filter for directories only
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return (dirs)
}

func prompt(reader *bufio.Reader, format string, a ...interface{}) string {
	// show prompt
	fmt.Printf(format, a...)

	// get response
	response, err := reader.ReadString('\n')
	checkError(err)

	return (response)
}

func main() {
	// read prompts from stdin
	stdin := bufio.NewReader(os.Stdin)

	// usage info
	fmt.Println("This program will read CSV files from a directory, add comma to EOL, and put new CSV file in current directory.\n")

	// get list of directories in current directory
	dirs := getDirs(".")

	if dirs == nil {
		fmt.Println("There are no directories to process")
		prompt(stdin, "Press <enter> to exit")
		os.Exit(1)
	}

	// loop thru each directory
	for _, dir := range dirs {
		// ask if this is a directory to process
		response := prompt(stdin, "Process directory %s [y/n]? ", dir)

		// if response beings with a y, then process this directory
		if strings.HasPrefix(response, "y") {
			matches, err := filepath.Glob(dir + string(os.PathSeparator) + "*.csv")
			checkError(err)

			// loop thru each file
			for _, fileName := range matches {
				processFile(fileName)
			}
		}
	}

}
