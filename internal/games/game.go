package games

import "errors"

var allWords []string = []string{"anpa", "ante", "awen", "esun", "insa", "jaki", "jelo", "kala", "kama", "kasi", "kili", "kule", "kute", "lape", "laso", "lawa", "lete", "lili", "lipu", "loje", "luka", "lupa", "mama", "mani", "moku", "moli", "musi", "mute", "nasa", "nena", "nimi", "noka", "olin", "open", "pali", "pana", "pini", "pipi", "poka", "poki", "pona", "sama", "seli", "selo", "seme", "sewi", "sike", "sina", "sona", "suli", "suno", "supa", "suwi", "taso", "tawa", "telo", "toki", "tomo", "unpa", "walo", "waso", "wawa", "weka", "wile", "leko", "meli", "mije", "soko"}

type Game struct {
	possibleWords []string
}

func NewGame() Game {
	return Game{
		possibleWords: allWords,
	}
}

func (g Game) MakeGuess() (string, float64, error) {
	minWordsAfterBestGuess := float64(len(allWords))
	bestGuess := ""
	for _, guess := range allWords {
		wordsAfterGuessAmounts := []int{}
		for _, word := range g.possibleWords {
			information := GetInformation(guess, word)
			wordsAfterGuess := 0
			for _, wordAfterGuess := range g.possibleWords {
				if wordAfterGuess != guess && information.Matches(wordAfterGuess) {
					wordsAfterGuess += 1
				}
			}
			wordsAfterGuessAmounts = append(wordsAfterGuessAmounts, wordsAfterGuess)
		}
		average := avg(wordsAfterGuessAmounts)
		if average < minWordsAfterBestGuess {
			minWordsAfterBestGuess = average
			bestGuess = guess
		}
	}
	if bestGuess == "" {
		return "", 0, errors.New("couldn't find a valid guess")
	}
	return bestGuess, minWordsAfterBestGuess, nil
}

func (g *Game) Information(information Information) int {
	newPossibleWords := []string{}
	for _, possibleWord := range g.possibleWords {
		if information.Matches(possibleWord) {
			newPossibleWords = append(newPossibleWords, possibleWord)
		}
	}
	g.possibleWords = newPossibleWords
	return len(newPossibleWords)
}

func avg(arr []int) float64 {
	sum := 0
	for _, i := range arr {
		sum += i
	}
	return float64(sum) / float64(len(arr))
}
