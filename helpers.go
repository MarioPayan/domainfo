package main

import (
	"log"
	"regexp"
)

var Ready = "Ready"
var READY = "READY"
var IN_PROGRESS = "IN_PROGRESS"
var ERROR = "ERROR"
var DNS = "DNS"

func manageError(err error, message string) {
	if err != nil {
		log.Fatalf(message+"\n ERROR: %s", err)
	}
}

func getRegexMatch(regex []string, text string, itemName string, sourceName string) string {
	var item string
	for _, iRegex := range regex {
		regularExpression, err := regexp.Compile("(?i)" + iRegex)
		manageError(err, "regex match failed")
		subMatchs := regularExpression.FindStringSubmatch(text)
		if len(subMatchs) >= 2 && subMatchs[1] != "" {
			item = subMatchs[1]
		}
	}
	return item
}
