package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type QuizItem struct {
	question string
	answer   string
}

func readFile(filename string) (quizItems []QuizItem) {

	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Could not open file: " + filename)
	}

	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		quizItem := QuizItem{question: record[0], answer: record[1]}
		quizItems = append(quizItems, quizItem)

	}

	return
}

func runQuiz(quizItems []QuizItem, timeLimit int) {
	numQuestions := len(quizItems)
	numCorrect := 0

	fmt.Print("Press Enter to begin ...")
	fmt.Scanln()

QuestionLoop:
	for i, item := range quizItems {

		fmt.Printf("Problem %d: %s = ", i+1, item.question)

		var answer string
		answerChannel := make(chan string)
		go func() {
			fmt.Scanln(&answer)
			answerChannel <- answer
		}()

		select {
		case <-answerChannel:
			// verify the anwer was correct
			if strings.EqualFold(strings.TrimSpace(answer), item.answer) {
				numCorrect += 1
			}
			break
		case <-time.After(time.Duration(timeLimit * int(time.Second))):
			fmt.Println("\nTime Limit Exceeded")
			break QuestionLoop
		}

	}

	fmt.Printf("You scored %d out of %d", numCorrect, numQuestions)

}

func main() {
	// commandLineArgs = os.Args[1:]
	var questionsFile string
	var shuffle bool
	var timeLimit int
	flag.StringVar(&questionsFile, "questions", "problems.csv", "The questions file")
	flag.BoolVar(&shuffle, "shuffle", false, "Shuffle the question order")
	flag.IntVar(&timeLimit, "timelimit", 30, "Time limit per question")

	flag.Parse()

	quizItems := readFile(questionsFile)

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizItems), func(i, j int) { quizItems[i], quizItems[j] = quizItems[j], quizItems[i] })
	}

	runQuiz(quizItems, timeLimit)
}
