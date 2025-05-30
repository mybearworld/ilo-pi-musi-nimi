package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/mybearworld/ilo-pi-musi-nimi/internal/games"
)

type row struct {
	guess            string
	averageWordsLeft float64
	realWordsLeft    int
	information      games.Information
}

var starter = flag.String("starter", "", "Pick a starter word instead of using the best one. Use \"random\" to choose one at random.")

func main() {
	flag.Parse()
	rl, err := readline.New("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed initializing readline: %v\n", err)
		os.Exit(1)
	}
	game := games.NewGame()
	rows := []row{}
	for nthGuess := 0; ; nthGuess++ {
		var (
			guess            string
			averageWordsLeft float64
		)
		if nthGuess == 0 && *starter != "" {
			guess = *starter
			if guess == "random" {
				guess = games.AllWords[rand.N(len(games.AllWords)-1)]
			}
			averageWordsLeft = game.ScoreGuess(guess)
		} else {
			guess, averageWordsLeft, err = game.MakeGuess()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed getting the next guess: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("I guess %s.\n", guess)
		information, err := inputInformation(rl, guess)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed getting input: %v\n", err)
			os.Exit(1)
		}
		realWordsLeft := game.Information(information)
		rows = append(rows, row{
			guess: guess, averageWordsLeft: averageWordsLeft, realWordsLeft: realWordsLeft, information: information,
		})
		if information.Success() {
			break
		}
	}
	fmt.Println()
	for _, row := range rows {
		realWordsLeftString := strconv.Itoa(row.realWordsLeft)
		if row.information.Success() {
			realWordsLeftString = "🥳"
		}
		fmt.Printf("%s %s %f %s\n", emojify(row.information), row.guess, row.averageWordsLeft, realWordsLeftString)
	}
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
