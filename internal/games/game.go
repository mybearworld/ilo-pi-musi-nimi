package games

import (
	"errors"
	"math/rand/v2"
)

var AllWords []string = []string{"anpa", "ante", "awen", "esun", "insa", "jaki", "jelo", "kala", "kama", "kasi", "kili", "kule", "kute", "lape", "laso", "lawa", "lete", "lili", "lipu", "loje", "luka", "lupa", "mama", "mani", "moku", "moli", "musi", "mute", "nasa", "nena", "nimi", "noka", "olin", "open", "pali", "pana", "pini", "pipi", "poka", "poki", "pona", "sama", "seli", "selo", "seme", "sewi", "sike", "sina", "sona", "suli", "suno", "supa", "suwi", "taso", "tawa", "telo", "toki", "tomo", "unpa", "walo", "waso", "wawa", "weka", "wile", "leko", "meli", "mije", "soko"}

type Strategy string

const (
	MinWords Strategy = "minwords"
	MaxWords Strategy = "maxwords"
	Random   Strategy = "random"
)

func ToStrategy(s string) *Strategy {
	strategy := Strategy(s)
	switch strategy {
	case MinWords, MaxWords, Random:
		return &strategy
	default:
		return nil
	}
}

type Game struct {
	possibleWords []string
	Strategy      Strategy
	Hard          bool
}

func NewGame(strategy Strategy, hard bool) Game {
	return Game{
		possibleWords: AllWords,
		Strategy:      strategy,
		Hard:          hard,
	}
}

func (g Game) MakeGuess() (string, float64, error) {
	switch g.Strategy {
	case MinWords, MaxWords:
		return g.makeGuessByWordScore()
	case Random:
		return g.makeRandomGuess()
	default:
		return "", 0, errors.New("unknown strategy")
	}
}

func (g Game) makeGuessByWordScore() (string, float64, error) {
	bestGuessScore := float64(len(AllWords))
	if g.Strategy == MaxWords {
		bestGuessScore = -1
	}
	bestGuess := ""
	for _, guess := range g.pool() {
		guessScore := g.ScoreGuess(guess)
		if (g.Strategy == MinWords && guessScore < bestGuessScore) || (g.Strategy == MaxWords && guessScore > bestGuessScore) {
			bestGuessScore = guessScore
			bestGuess = guess
		}
	}
	if bestGuess == "" {
		return "", 0, errors.New("couldn't find a valid guess")
	}
	return bestGuess, bestGuessScore, nil
}

func (g Game) makeRandomGuess() (string, float64, error) {
	pool := g.pool()
	guess := pool[rand.N(len(pool))]
	return guess, g.ScoreGuess(guess), nil
}

func (g Game) pool() []string {
	guessPool := AllWords
	if g.Hard {
		guessPool = g.possibleWords
	}
	return guessPool
}

func (g Game) ScoreGuess(guess string) float64 {
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
	return average
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
