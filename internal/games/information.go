package games

import (
	"fmt"
)

type Information struct {
	GreenLetters    []LetterWithIndex
	NonGreenLetters []NonGreenLetter
}

func (i Information) String() string {
	return fmt.Sprintf("{GreenLetters: %s, NonGreenLetters: %s}", i.GreenLetters, i.NonGreenLetters)
}

func (i Information) Success() bool {
	return len(i.NonGreenLetters) == 0
}

func (i Information) Matches(word string) bool {
	wordRunes := []rune(word)
	correctPositionIndices := map[int]struct{}{}
	for _, letter := range i.GreenLetters {
		if wordRunes[letter.Index] != letter.Letter {
			return false
		}
		correctPositionIndices[letter.Index] = struct{}{}
	}
	checkedYellows := map[LetterWithIndex]struct{}{}
	yellowIndices := map[int]struct{}{}
	for _, letter := range i.NonGreenLetters {
		if letter.IsYellow {
			foundMatch := false
			for wordLetterIndex, wordLetter := range word {
				if wordLetter == letter.Letter {
					if wordLetterIndex == letter.Index {
						return false
					}
					if _, ok := checkedYellows[LetterWithIndex{Letter: wordLetter, Index: wordLetterIndex}]; ok {
						continue
					}
					checkedYellows[LetterWithIndex{Letter: wordLetter, Index: wordLetterIndex}] = struct{}{}
					yellowIndices[wordLetterIndex] = struct{}{}
					foundMatch = true
				}
			}
			if !foundMatch {
				return false
			}
		} else {
			for index, wordLetter := range wordRunes {
				if wordLetter == letter.Letter {
					if _, ok := correctPositionIndices[index]; ok {
						continue
					}
					if _, ok := yellowIndices[index]; ok {
						continue
					}
					return false
				}
			}
		}
	}
	return true
}

type LetterWithIndex struct {
	Letter rune
	Index  int
}

func (gl LetterWithIndex) String() string {
	return fmt.Sprintf("%c:%d", gl.Letter, gl.Index)
}

type NonGreenLetter struct {
	Letter   rune
	Index    int
	IsYellow bool
}

func (ngl NonGreenLetter) String() string {
	colorName := "Gray"
	if ngl.IsYellow {
		colorName = "Yellow"
	}
	return fmt.Sprintf("%s:%c:%d", colorName, ngl.Letter, ngl.Index)
}

func GetInformation(guess string, correctWord string) Information {
	correctWordRunes := []rune(correctWord)
	information := Information{}
	correctPositionIndices := map[int]struct{}{}
	for index, letter := range guess {
		if correctWordRunes[index] == letter {
			information.GreenLetters = append(information.GreenLetters, LetterWithIndex{
				Letter: letter, Index: index,
			})
			correctPositionIndices[index] = struct{}{}
		}
	}
	letterPositions := map[rune]int{}
	for index, letter := range correctWord {
		if _, ok := correctPositionIndices[index]; !ok {
			letterPositions[letter] += 1
		}
	}
	for index, letter := range guess {
		if _, ok := correctPositionIndices[index]; ok {
			continue
		}
		if letterPositions[letter] == 0 {
			information.NonGreenLetters = append(information.NonGreenLetters, NonGreenLetter{
				Letter: letter, Index: index, IsYellow: false,
			})
		} else {
			letterPositions[letter]--
			information.NonGreenLetters = append(information.NonGreenLetters, NonGreenLetter{
				Letter: letter, Index: index, IsYellow: true,
			})
		}
	}
	return information
}
