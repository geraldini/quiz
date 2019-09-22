package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question, CorrectAnswer, UserAnswer string
}

func (problem *Problem) Ask(index int) bool {
	fmt.Printf("Question #%d: %s\nYour answer: ", index, problem.Question)
	userAnswer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	problem.UserAnswer = strings.TrimSpace(userAnswer)
	return strings.Compare(problem.CorrectAnswer, problem.UserAnswer) == 0
}

type Quiz struct {
	CorrectAnswers, TimeLimit int
	Problems                  []Problem
}

func (quiz *Quiz) Execute() {
	fmt.Println("Welcome to today's quiz! Please press ENTER to start.")
	defer quiz.PrintSummary()
	bufio.NewReader(os.Stdin).ReadString('\n')
	quizChannel := make(chan string)
	timerChannel := make(chan string)
	go quiz.AskQuestions(quizChannel)
	go quiz.TimeQuiz(timerChannel)
	select {
	case <-timerChannel:
		return
	case <-quizChannel:
		return
	}
}

func (quiz *Quiz) AskQuestions(channel chan string) {
	for index, problem := range quiz.Problems {
		answerIsCorrect := problem.Ask(index + 1)
		if answerIsCorrect {
			quiz.CorrectAnswers++
		}
	}
	fmt.Println("Quiz Completed!")
	channel <- "Completed"
}

func (quiz *Quiz) TimeQuiz(channel chan string) {
	timer := time.NewTimer(time.Duration(quiz.TimeLimit) * time.Second)
	<-timer.C
	fmt.Println("Time's up!")
	channel <- "Time's Up!"
}

func (quiz *Quiz) PrintSummary() {
	fmt.Printf("Total Questions: %d\nCorrect Answers: %d\nThank you for playing!", len(quiz.Problems), quiz.CorrectAnswers)
}

func (quiz *Quiz) LoadQuestions(filePath string) {
	csvFile, err := os.Open(filePath)
	defer csvFile.Close()
	if err != nil {
		log.Fatal("Couldn't open CSV file" + filePath)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Couldn't parse CSV file" + filePath)
	}
	quiz.RandomizeProblems(lines)
}

func (quiz *Quiz) RandomizeProblems(lines [][]string) {
	nProblems := len(lines)
	quiz.Problems = make([]Problem, nProblems)
	rand.Seed(time.Now().UnixNano())
	randomOrder := rand.Perm(nProblems)
	for originalIndex, randomIndex := range randomOrder {
		quiz.Problems[randomIndex] = Problem{
			Question:      strings.TrimSpace(lines[originalIndex][0]),
			CorrectAnswer: strings.TrimSpace(lines[originalIndex][1]),
		}
	}
}

func main() {
	filePath := flag.String("file-path", "quiz1.csv", "Path to the CSV file with the questions")
	timeLimit := flag.Int("time-limit", 30, "Time allowed to complete the quiz")
	flag.Parse()
	quiz := Quiz{TimeLimit: *timeLimit}
	quiz.LoadQuestions(*filePath)
	quiz.Execute()
}
