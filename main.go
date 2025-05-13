package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// calculate parses and evaluates a mathematical expression string.
// It removes all spaces from the input and delegates evaluation to evalExpr.
// Returns the computed result or an error if the expression is invalid.
func calculate(input string) (float64, error) {
	expr := strings.ReplaceAll(input, " ", "")
	result, err := evalExpr(expr)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// evalExpr evaluates a mathematical expression string supporting +, -, *, /, and parentheses.
// It uses recursive descent parsing to handle operator precedence and parentheses.
// Returns the computed result or an error if the expression is invalid.
func evalExpr(expr string) (float64, error) {
	var parse func() (float64, error)
	tokens := []rune(expr)
	pos := 0

	// parseFactor parses numbers and parenthesized sub-expressions.
	var parseFactor func() (float64, error)
	parseFactor = func() (float64, error) {
		// Skip whitespace (shouldn't be any)
		for pos < len(tokens) && tokens[pos] == ' ' {
			pos++
		}
		if pos < len(tokens) && tokens[pos] == '(' {
			pos++
			val, err := parse()
			if err != nil {
				return 0, err
			}
			if pos >= len(tokens) || tokens[pos] != ')' {
				return 0, fmt.Errorf("missing closing parenthesis")
			}
			pos++
			return val, nil
		}
		start := pos
		dotSeen := false
		for pos < len(tokens) && (tokens[pos] >= '0' && tokens[pos] <= '9' || tokens[pos] == '.') {
			if tokens[pos] == '.' {
				if dotSeen {
					return 0, fmt.Errorf("invalid number format")
				}
				dotSeen = true
			}
			pos++
		}
		if start == pos {
			return 0, fmt.Errorf("expected number at position %d", pos)
		}
		num, err := strconv.ParseFloat(string(tokens[start:pos]), 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	}

	// parseTerm parses multiplication and division operations.
	var parseTerm func() (float64, error)
	parseTerm = func() (float64, error) {
		val, err := parseFactor()
		if err != nil {
			return 0, err
		}
		for pos < len(tokens) {
			if tokens[pos] == '*' || tokens[pos] == '/' {
				op := tokens[pos]
				pos++
				nextVal, err := parseFactor()
				if err != nil {
					return 0, err
				}
				if op == '*' {
					val *= nextVal
				} else {
					if nextVal == 0 {
						return 0, fmt.Errorf("division by zero")
					}
					val /= nextVal
				}
			} else {
				break
			}
		}
		return val, nil
	}

	// parse parses addition and subtraction operations.
	parse = func() (float64, error) {
		val, err := parseTerm()
		if err != nil {
			return 0, err
		}
		for pos < len(tokens) {
			if tokens[pos] == '+' || tokens[pos] == '-' {
				op := tokens[pos]
				pos++
				nextVal, err := parseTerm()
				if err != nil {
					return 0, err
				}
				if op == '+' {
					val += nextVal
				} else {
					val -= nextVal
				}
			} else {
				break
			}
		}
		return val, nil
	}

	result, err := parse()
	if err != nil {
		return 0, err
	}
	if pos != len(tokens) {
		return 0, fmt.Errorf("unexpected character at position %d", pos)
	}
	return result, nil
}

// main is the entry point of the calculator program.
// It reads user input from stdin, evaluates mathematical expressions, and prints the result.
// The program exits when the user enters 'exit'.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter calculation (<number> <operator> <number>), or 'exit' to quit:")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			break
		}
		result, err := calculate(line)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	}
}
