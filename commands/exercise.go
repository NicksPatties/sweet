package main

import "fmt"

func RunExercise(lang string, topic string) {
	if lang == "" {
		lang = "any"
	}
	if topic == "" {
		topic = "any"
	}
	fmt.Printf("Run exercise\n")
	fmt.Printf("  lang:  %s\n", lang)
	fmt.Printf("  topic: %s\n", topic)
}
