package pkg

import (
	"log"
)

func Handler(err error) {
	if err != nil {
		log.Println(err)
	}
	return
}
