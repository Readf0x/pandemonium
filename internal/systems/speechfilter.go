package systems

import goaway "github.com/TwiN/go-away"

var (
	dict = []string{
		"coon",
		"fag",
		"nazi",
		"nigger",
		"nigga",
		"niggu",
		"queer",
		"rape",
		"rapist",
		"tranny",
		"troon",
	}

	SpeechFilter = goaway.NewProfanityDetector().WithCustomDictionary(dict, goaway.DefaultFalsePositives, goaway.DefaultFalseNegatives)
)

