package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var vowels = "aeiou"
var consonants = "bcdfghjklmnpqrstvwxyz"
var freqMap = map[byte][]byte{}

var startSyllables, middleSyllables, endSyllables []string

func initFreqMap() {
	for i := 0; i < len(consonants); i++ {
		freqMap[consonants[i]] = []byte(vowels)
	}
	for i := 0; i < len(vowels); i++ {
		freqMap[vowels[i]] = []byte(consonants)
	}
}

func loadSyllables(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var jsonData struct {
		Start  []string `json:"start"`
		Middle []string `json:"middle"`
		End    []string `json:"end"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	startSyllables = jsonData.Start
	middleSyllables = jsonData.Middle
	endSyllables = jsonData.End
	return nil
}

func selectSyllable(prev string, candidates []string, maxLength int) string {
	if len(prev) == 0 {
		return candidates[rand.Intn(len(candidates))]
	}
	last := prev[len(prev)-1]
	filtered := []string{}

	for _, s := range candidates {
		if len(s) > maxLength {
			continue
		}
		if chars, ok := freqMap[last]; ok && strings.Contains(string(chars), string(s[0])) {
			filtered = append(filtered, s)
		}
	}

	if len(filtered) == 0 {
		for _, s := range candidates {
			if len(s) <= maxLength {
				filtered = append(filtered, s)
			}
		}
	}

	if len(filtered) == 0 {
		return ""
	}

	idx := rand.Intn(len(filtered))
	if rand.Float32() < 0.3 {
		idx = (idx + 1) % len(filtered)
	}

	if len(filtered[idx]) == 2 && strings.Contains(consonants, string(filtered[idx][0])) && rand.Float32() < 0.1 {
		filtered[idx] = string(filtered[idx][0]) + filtered[idx]
	}

	if len(filtered[idx]) < 3 && rand.Float32() < 0.15 {
		filtered[idx] += string(vowels[rand.Intn(len(vowels))])
	}

	return filtered[idx]
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func generateNickname() string {
	targetLen := 8
	nick := startSyllables[rand.Intn(len(startSyllables))]

	for len(nick) < targetLen {
		remain := targetLen - len(nick)
		var s string
		if remain <= 2 {
			s = selectSyllable(nick, endSyllables, remain)
		} else {
			s = selectSyllable(nick, middleSyllables, remain)
		}
		if s == "" {
			break
		}
		if len(nick)+len(s) > targetLen {
			s = s[:targetLen-len(nick)]
		}
		nick += s
	}

	return capitalize(nick)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initFreqMap()

	if err := loadSyllables("syllables.json"); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(generateNickname())
	}
}