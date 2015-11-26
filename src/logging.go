package scarecrow

import "fmt"

// Logging-related utility functions.
// TODO: better logging.

func (self *Scarecrow) Log(message string, a ...interface{}) {
	if self.Debug {
		fmt.Printf("[DEBUG] "+message+"\n", a...)
	}
}

func (self *Scarecrow) Info(message string, a ...interface{}) {
	fmt.Printf("[INFO] "+message+"\n", a...)
}

func (self *Scarecrow) Warn(message string, a ...interface{}) {
	fmt.Printf("[WARN] "+message+"\n", a...)
}

func (self *Scarecrow) Error(message string, a ...interface{}) {
	fmt.Printf("[Error] "+message+"\n", a...)
}
