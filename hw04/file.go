package tinyini

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"strings"
)

// Determines whether the file path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Reads the file contents into memory
func readFile(filename string) (*os.File, error) {
	if !PathExists(filename) {
		fmt.Println("Error: ", ErrNoFile.Error())
		return nil, ErrNoFile
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", ErrOpenFile.Error())
		return nil, ErrOpenFile
	}
	return file, nil
}

// Load file information in memory
func loadFile(filename string) (*config, error) {
	file, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := new(config)
	c.filePath = filename
	c.info = make(map[string]map[string]string)

	var section, key, value string
	
	buf := bufio.NewReader(file)
	// Read file information line by line
	for {
		// Read a line of information
		line, err := buf.ReadString('\n')
		// Removes Spaces at the beginning and end of the string
		line = strings.TrimSpace(line)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", ErrReadFile.Error())
				return c, err
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
			// Ignore the null line
		case string(line[0]) == annotationFlag:
			// Ignore the annotation
		case line[0] == '[' && line[len(line)-1] == ']':
			section = line[1 : len(line)-1]
			c.info[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			key = strings.TrimSpace(line[0:i])
			value = strings.TrimSpace(line[i+1 : len(line)])

			// There is no section in the file
			if section == "" {
				section = "NULL"
				c.info[section] = make(map[string]string)
			}
			c.info[section][key] = value
		}
	}
	return c, nil
}