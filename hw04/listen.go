package tinyini

import "fmt"

func keepLoadListener(filename string) (*config, error) {
	c, err := loadFile(filename)
	if err != nil {
		return nil, err
	}
	c.display()
	// Read the file continuously
	for {
		newc, err := loadFile(filename)
		if err != nil {
			return nil, err
		}
		// When changes are found, reprint the information and save the new config
		if (!isSame(c, newc)) {
			fmt.Println("File has changed!")
			newc.display()
			c = newc
		}
	}
	return c, nil
}