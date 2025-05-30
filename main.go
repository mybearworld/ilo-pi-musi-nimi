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

var flagDictionary = flag.String("dictionary", "core,common", "The words to allow as the solution. A comma separated list of core, common, uncommon, obscure or sanbox, which correspond to their respective Linku categories.")
var flagGuessDictionary = flag.String("guessdictionary", "", "The words to allow for guessing. A comma separated list of core, common, uncommon, obscure or sanbox, which correspond to their respective Linku categories. (default is the regular word dictionary)")
var starter = flag.String("starter", "", "Pick a starter word instead of using the strategy. Use \"random\" to choose one at random.")
var word = flag.String("word", "", "Pick a word to solve for, instead of you needing to enter the inputs manually. Use \"random\" to choose one at random.")
var flagStrategy = flag.String("strategy", "minwords", "The strategy to guess with. Your options are:\n- minwords: Chooses the word that'll leave the fewest amount of words.\n- maxwords: Chooses the word that'll leave the highest amount of words.\n- random: Picks words at random.\n")
var hard = flag.Bool("hard", false, "Require using previous hints in subsequent guesses.")
var guesses = flag.Int("guesses", 6, "The amount of guesses to allow. Set to 0 for unlimited guesses. (This might loop infinitely for some configurations.)")
var strategy games.Strategy
var dictionary = []string{}
var guessDictionary = []string{}

func parseFlags() string {
	flag.Parse()
	if msg := parseDictionary(*flagDictionary, &dictionary); msg != "" {
		return msg
	}
	if *flagGuessDictionary == "" {
		guessDictionary = dictionary
	} else if msg := parseDictionary(*flagGuessDictionary, &guessDictionary); msg != "" {
		return msg
	}
	if *starter == "random" {
		*starter = dictionary[rand.N(len(dictionary))]
	} else if *starter != "" && len(*starter) != 4 {
		return "Starter must be four characters long"
	}
	if *word == "random" {
		*word = dictionary[rand.N(len(dictionary))]
	} else if *word != "" && !slices.Contains(dictionary, *word) {
		return "Word must be a valid 4-letter toki pona word."
	}
	parsedStrategy := games.ToStrategy(*flagStrategy)
	if parsedStrategy == nil {
		return "Invalid strategy: " + *flagStrategy
	} else {
		strategy = *parsedStrategy
	}
	return ""
}

func parseDictionary(dictionaryString string, dictionary *[]string) string {
	wordCategories := strings.Split(dictionaryString, ",")
	for _, category := range wordCategories {
		switch category {
		case "core":
			*dictionary = append(*dictionary, games.CoreWords...)
		case "common":
			*dictionary = append(*dictionary, games.CommonWords...)
		case "uncommon":
			*dictionary = append(*dictionary, games.UncommonWords...)
		case "obscure":
			*dictionary = append(*dictionary, games.ObscureWords...)
		case "sandbox":
			*dictionary = append(*dictionary, games.SandboxWords...)
		case " ", "":
		default:
			return "Invalid word category: " + category
		}
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
	game := games.NewGame(dictionary, guessDictionary, strategy, *hard)
	rows := []row{}
	for nthGuess := 0; *guesses == 0 || nthGuess < *guesses; nthGuess++ {
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
		guessScore, err := game.Information(information)
		if err != nil {
			return "Failed processing information: " + err.Error()
		}
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
		if len(line) != 4 {
			fmt.Println("This needs to be four characters long.")
			continue
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
