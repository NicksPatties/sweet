package exercise

import "fmt"

// Runs the exercise, which is a bubbletea application.
// Returns a status code.
func Run(lang string, topic string, filename string) int {
	if lang == "" {
		lang = "any"
	}
	if topic == "" {
		topic = "any"
	}
	if filename == "" {
		filename = "random"
	}
	fmt.Printf("Run exercise\n")
	fmt.Printf("  lang:  %s\n", lang)
	fmt.Printf("  topic: %s\n", topic)
	fmt.Printf("  file: %s\n", filename)
	return 0
}
