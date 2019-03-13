package shared

import (
	"fmt"
	"log"
)

var Verbose = false

func VerbosePrintln(s string) {
	if Verbose {
		fmt.Println(s)
	}
}

func CheckErrorAndFatal(description string, err error) {
	if err != nil {
		log.Fatalln(description+":", err)
	}
}

func CheckErrorAndCleanAndFatal(description string, err error, cleanUp func()) {
	if err != nil {
		cleanUp()
		log.Fatalln(description+":", err)
	}
}
