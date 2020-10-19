package tinyini

import "fmt"

// Defines the structure of configuration information
type config struct {
	filePath string
	info map[string]map[string]string
}

// Determine if two configs are identical
func isSame(c1 *config, c2 *config) bool {
	for s1, m1 := range c1.info {
		for k1, v1 := range m1 {
			_, ok := c2.info[s1]
			if !ok {
				return false
			}
			v2, ok := c2.info[s1][k1]
			if !ok {
				return false
			}
			if v2 != v1 {
				return false
			}
		}
	}
	return true
}

// Print the content information for the config
func (c *config) display() {
	fmt.Println("***********************************")
	for s, m := range c.info {
		fmt.Println("S: ", s)
		for k, v := range m {
			fmt.Println("  K: " + k)
			fmt.Println("      V: " + v)
		}
		fmt.Println(" ")
	}
	fmt.Println("***********************************")
}