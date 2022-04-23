package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Input      string
	Mistake    int
	Location   []int
	WrongWords []string
}

type Request struct {
	Sentence string
}

func checkSpelling(sentence string) *Response {

	var isIn bool

	response := new(Response)
	sentence = strings.TrimSpace(sentence)

	words := strings.Split(sentence, " ")

	var wg sync.WaitGroup
	wg.Add(len(words))
	var mu sync.Mutex
	for i, word := range words {
		go wordIsIn(i, word, isIn, response, &wg, &mu)
	}
	wg.Wait()
	sort.Ints(response.Location)
	response.Input = sentence
	return response
}

func wordIsIn(i int, word string, isIn bool, response *Response, wg *sync.WaitGroup, mu *sync.Mutex) {
	mu.Lock()
	word = strings.ToLower(word)

	if isANumber(word) {
		word = "a"
	}
	if isSpecial(word) {
		word = "a"
	}
	if isMix(word) {
		word = "aabababba"
	}
	if isMixWithLetter(word) {
		word = "ababababababab"
	}
	for _, v := range word {
		if v < 97 || v > 122 {
			word = strings.Replace(word, string(v), "", 1)
		}
	}

	beginningLetter := string(word[0])

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fileName := fmt.Sprintf((pwd + "/dictionary/%s.txt"), beginningLetter)

	f, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if line == word {
			isIn = true
		}
	}

	if !isIn {
		response.Mistake += 1
		response.Location = append(response.Location, i)
		response.Location = sort.IntSlice(response.Location)
		response.WrongWords = append(response.WrongWords, word)

	}

	if isIn {
		isIn = false
	}

	if response.WrongWords == nil {
		response.WrongWords = []string{}
	}

	if response.Location == nil {
		response.Location = []int{}
	}

	f.Close()
	mu.Unlock()
	wg.Done()
}

func getWord(ctx *fiber.Ctx) error {

	body := new(Request)

	err := ctx.BodyParser(body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	sentence := body.Sentence

	response := checkSpelling(sentence)

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func isANumber(word string) bool {
	_, err := strconv.Atoi(word)
	return err == nil
}

func isSpecial(word string) bool {
	if len(word) > 1 {
		return false
	}
	for _, v := range word {
		if v > 122 || v < 97 {
			return true
		}
	}
	return false
}

func isMix(word string) bool {
	amount := 0
	for _, v := range word {
		if v < 97 || v > 122 {
			amount++
		}
	}
	return amount == len(word)
}

func isMixWithLetter(word string) bool {
	amount := 0
	hasSpecial := false
	for _, v := range word {
		if v < 97 || v > 122 {
			amount++
		} else if v <= 122 && v >= 97 {
			hasSpecial = true
		}
	}
	return hasSpecial == true && amount >= 2
}

// 1.46ms per word
