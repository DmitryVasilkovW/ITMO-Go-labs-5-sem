package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	prefixForBaseCommands       = "b_"
	keyCharacterToCreateNewWord = ":"
	errorMessageForUnderflow    = "stack underflow"
)

type Evaluator struct {
	stack []int
	words map[string][]string
	base  map[string]func(*Evaluator) error
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		stack: make([]int, 0),
		words: make(map[string][]string),
		base: map[string]func(*Evaluator) error{
			"+":    plus,
			"-":    minus,
			"*":    multiply,
			"/":    divide,
			"dup":  dup,
			"drop": drop,
			"swap": swap,
			"over": over,
		},
	}
	return e
}

func (e *Evaluator) Process(row string) ([]int, error) {
	tokens := strings.Fields(row)
	var err error
	for i := 0; i < len(tokens); i++ {
		token := strings.ToLower(tokens[i])
		if token == keyCharacterToCreateNewWord {
			err = createUserWord(e, tokens, &i)
		} else {
			err = e.executeWord(token)
		}
	}

	if err != nil {
		return nil, err
	}

	return e.stack, nil
}

func createUserWord(e *Evaluator, tokens []string, indexOfStartOfToken *int) error {
	if *indexOfStartOfToken+1 >= len(tokens) {
		return errors.New("invalid definition")
	}

	word := strings.ToLower(tokens[*indexOfStartOfToken+1])
	if _, err := strconv.Atoi(word); err == nil {
		return fmt.Errorf("cannot redefine number: %s", word)
	}

	endIndex := findIndexOfEndOfToken(*indexOfStartOfToken, tokens)
	if endIndex == -1 {
		return errors.New("missing ';' in definition")
	}

	e.words[word] = createUserWordBody(e, tokens, *indexOfStartOfToken, endIndex)

	*indexOfStartOfToken = endIndex

	return nil
}

func createUserWordBody(e *Evaluator, tokens []string, indexOfStartOfToken, indexOfEndOfToken int) []string {
	var resultSubWords []string
	for _, subWord := range tokens[indexOfStartOfToken+2 : indexOfEndOfToken] {
		subWord = strings.ToLower(subWord)
		if _, exists := e.words[subWord]; exists {
			resultSubWords = append(resultSubWords, e.words[subWord]...)
		} else if _, err := strconv.Atoi(subWord); err == nil {
			resultSubWords = append(resultSubWords, subWord)
		} else {
			resultSubWords = append(resultSubWords, "b_"+subWord)
		}
	}

	return resultSubWords
}

func findIndexOfEndOfToken(indexOfStartOfToken int, tokens []string) int {
	endIndex := -1
	for j := indexOfStartOfToken + 2; j < len(tokens); j++ {
		if tokens[j] == ";" {
			endIndex = j
			break
		}
	}

	return endIndex
}

func (e *Evaluator) executeWord(token string) error {
	token = strings.ToLower(token)

	if def, isWord := e.words[token]; isWord {
		return executeUserWord(e, def)
	}

	if operation, exists := e.base[token]; exists {
		return operation(e)
	}

	if tryToAddInt(e, token) == nil {
		return nil
	}

	return fmt.Errorf("undefined word: %s", token)
}

func executeUserWord(e *Evaluator, def []string) error {
	for _, subToken := range def {

		if tryToAddInt(e, subToken) == nil {
			continue
		}

		err := executeSubWord(e, subToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeSubWord(e *Evaluator, subToken string) error {
	baseWord := subToken[2:]
	var err error
	if subToken[:2] == prefixForBaseCommands {
		if operation, exists := e.base[baseWord]; exists {
			err = operation(e)
		}
	} else {
		err = e.executeWord(subToken)
	}

	return err
}

func tryToAddInt(e *Evaluator, subToken string) error {
	if value, numErr := strconv.Atoi(subToken); numErr == nil {
		e.stack = append(e.stack, value)
		return nil
	} else {
		return numErr
	}
}

func plus(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	a, b := popTwoElements(&e.stack)
	e.stack = append(e.stack, a+b)
	return nil
}

func minus(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	a, b := popTwoElements(&e.stack)
	e.stack = append(e.stack, a-b)
	return nil
}

func multiply(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	a, b := popTwoElements(&e.stack)
	e.stack = append(e.stack, a*b)
	return nil
}

func divide(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	a, b := popTwoElements(&e.stack)
	if b == 0 {
		return errors.New("division by zero")
	}
	e.stack = append(e.stack, a/b)
	return nil
}

func dup(e *Evaluator) error {
	if len(e.stack) < 1 {
		return errors.New(errorMessageForUnderflow)
	}
	top := e.stack[len(e.stack)-1]
	e.stack = append(e.stack, top)
	return nil
}

func drop(e *Evaluator) error {
	if len(e.stack) < 1 {
		return errors.New(errorMessageForUnderflow)
	}
	e.stack = e.stack[:len(e.stack)-1]
	return nil
}

func swap(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	a, b := popTwoElements(&e.stack)
	e.stack = append(e.stack, b, a)
	return nil
}

func popTwoElements(stack *[]int) (int, int) {
	s := *stack
	a := s[len(s)-2]
	b := s[len(s)-1]
	*stack = s[:len(s)-2]
	return a, b
}

func over(e *Evaluator) error {
	if len(e.stack) < 2 {
		return errors.New(errorMessageForUnderflow)
	}
	second := e.stack[len(e.stack)-2]
	e.stack = append(e.stack, second)
	return nil
}
