package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cheggaaa/pb"
)

// a compound word
type Compound struct {
	word          string
	originalWords []string
}

func (c *Compound) String() string {
	return fmt.Sprintf("%s {%v}", c.word, c.originalWords)
}

// a list of compound words, and the index of the longest one
type CompoundList struct {
	longestIdx    int
	compoundWords []Compound
}

func (r *CompoundList) GetLongest() Compound {
	return r.compoundWords[r.longestIdx]
}

// from a given list of words, find all the words that are a compound of at least
// least two other words in the list, and determine which one is the longest.
func FindCompoundWords(words []string, includeProgressBar bool) (result *CompoundList) {
	result = &CompoundList{0, make([]Compound, 0, 128)}

	// this slice stores words that have appeared as sub-words before
	likelySubWords := make([]string, 0, 128)

	var bar *pb.ProgressBar

	if includeProgressBar {
		bar = pb.StartNew(len(words))
	}

	for _, currentWord := range words {
		if len(currentWord) == 0 {
			continue
		}

		subWords := make([]string, 0, 2)
		w := currentWord

		// first check the list of words that have already appeared as sub-words
		for _, subWord := range likelySubWords {
			if currentWord != subWord && strings.Contains(w, subWord) {
				w = strings.Replace(w, subWord, "", 1)
				subWords = append(subWords, subWord)
			}

			if len(subWords) >= 2 {
				break
			}
		}

		if len(subWords) < 2 {
			// fall back to searching the full list
			for _, subWord := range words {
				if currentWord != subWord && strings.Contains(w, subWord) {
					w = strings.Replace(w, subWord, "", 1)
					subWords = append(subWords, subWord)
					likelySubWords = append(likelySubWords, subWord)
				}

				if len(subWords) >= 2 {
					break
				}
			}
		}

		if len(subWords) >= 2 {
			c := Compound{currentWord, subWords}
			result.compoundWords = append(result.compoundWords, c)

			if len(currentWord) > len(result.GetLongest().word) {
				result.longestIdx = len(result.compoundWords) - 1
			}
		}

		if includeProgressBar {
			bar.Increment()
		}
	}

	if includeProgressBar {
		bar.FinishPrint(fmt.Sprintf("Processed %d words.", len(words)))
	}

	return result
}

// take a string and transform it into a list of words with whitespace trimmed
// and any empty elements removed
func wordsFromString(str string) (words []string) {
	words = strings.Split(str, "\n")

	toDelete := make([]int, 0, 8)

	for i, word := range words {
		word = strings.TrimSpace(word)
		if len(word) == 0 {
			toDelete = append(toDelete, i)
		}
	}

	for _, target := range toDelete {
		words = append(words[:target], words[target+1:]...)
	}

	return words
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "the file to parse")

	var includeProgressBar bool
	flag.BoolVar(&includeProgressBar, "progress", false, "show progress bar")

	flag.Parse()

	if len(inputFile) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	words := wordsFromString(string(data))
	compounds := FindCompoundWords(words, includeProgressBar)

	fmt.Println("Longest word was", compounds.GetLongest(), "with length", len(compounds.GetLongest().word))
}
