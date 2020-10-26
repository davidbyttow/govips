package vips

import "log"

func info(fmt string, values ...interface{}) {
	if len(values) > 0 {
		log.Printf(fmt, values...)
	} else {
		log.Print(fmt)
	}
}
