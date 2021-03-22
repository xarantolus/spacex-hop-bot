package jobs

import "log"

func logError(err error, location string) {
	if err != nil {
		log.Printf("[Error (%s)]: %s\n", location, err.Error())
	}
}
