package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/mybearworld/ilo-pi-musi-nimi/internal/games"
)

type row struct {
	guess               string
	projectedGuessScore float64
	guessScore          int
	information         games.Information
}

var starter = flag.String("starter", "", "Pick a starter word instead of using the best one. Use \"random\" to choose one at random.")
var word = flag.String("word", "", "Pick a word to solve for, instead of you needing to enter the inputs manually.")
var worst = flag.Bool("worst", false, "Always use the worst available guess instead of the best one. This option works best with -hard.")
var hard = flag.Bool("hard", false, "Require using previous hints in subsequent guesses.")

func parseFlags() string {
	flag.Parse()
	if *starter == "random" {
		*starter = games.AllWords[rand.N(len(games.AllWords))]
	} else if *starter != "" && len(*starter) != 4 {
		return "Starter must be four characters long"
	}
	if *word != "" && !slices.Contains(games.AllWords, *word) {
		return "Word must be a valid 4-letter toki pona word."
	}
	return ""
}

func main() {
	if msg := run(); msg != "" {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

func run() string {
	if msg := parseFlags(); msg != "" {
		return msg
	}
	rl, err := readline.New("")
	if err != nil {
		return "Failed initializing readline: " + err.Error()
	}
	game := games.NewGame(*worst, *hard)
	rows := []row{}
	for nthGuess := 0; nthGuess < 6; nthGuess++ {
		var (
			guess               string
			projectedGuessScore float64
		)
		if nthGuess == 0 && *starter != "" {
			guess = *starter
			projectedGuessScore = game.ScoreGuess(guess)
		} else {
			guess, projectedGuessScore, err = game.MakeGuess()
			if err != nil {
				return "Failed getting the next guess: " + err.Error()
			}
		}
		if *word == "" {
			fmt.Printf("I guess %s.\n", guess)
		}
		var information games.Information
		if *word != "" {
			information = games.GetInformation(guess, *word)
		} else {
			information, err = inputInformation(rl, guess)
			if err != nil {
				return "Failed getting input: " + err.Error()
			}
		}
		guessScore := game.Information(information)
		rows = append(rows, row{
			guess: guess, projectedGuessScore: projectedGuessScore, guessScore: guessScore, information: information,
		})
		if information.Success() {
			break
		}
	}
	if *word == "" {
		fmt.Println()
	}
	for _, row := range rows {
		realWordsLeftString := strconv.Itoa(row.guessScore)
		if row.information.Success() {
			realWordsLeftString = "🥳"
		}
		fmt.Printf("%s %s %f %s\n", emojify(row.information), row.guess, row.projectedGuessScore, realWordsLeftString)
	}
	return ""
}

func inputInformation(rl *readline.Instance, guess string) (games.Information, error) {
	guessRunes := []rune(guess)
	for {
		information := games.Information{}
		line, err := rl.Readline()
		if err != nil {
			return information, err
		}
		succeeded := true
		for index, letter := range line {
			switch letter {
			case 'l':
				information.GreenLetters = append(information.GreenLetters, games.LetterWithIndex{
					Letter: guessRunes[index], Index: index,
				})
			case 'j', 'p':
				information.NonGreenLetters = append(information.NonGreenLetters, games.NonGreenLetter{
					Letter: guessRunes[index], Index: index, IsYellow: letter == 'j',
				})
			default:
				succeeded = false
				fmt.Printf("I don't know what %c means.\n", letter)
			}
		}
		if succeeded {
			return information, nil
		}
	}
}

func emojify(information games.Information) string {
	emoji := map[int]string{}
	maxIndex := 0
	for _, letter := range information.GreenLetters {
		emoji[letter.Index] = "🟩"
		if letter.Index > maxIndex {
			maxIndex = letter.Index
		}
	}
	for _, letter := range information.NonGreenLetters {
		if letter.IsYellow {
			emoji[letter.Index] = "🟨"
		} else {
			emoji[letter.Index] = "🔳"
		}
		if letter.Index > maxIndex {
			maxIndex = letter.Index
		}
	}
	builder := strings.Builder{}
	for index := range maxIndex + 1 {
		builder.WriteString(emoji[index])
	}
	return builder.String()
}
