package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ExecuteBinOps(seq []string, pos int, sign string) (string, error) {
	// выполняет выражение и возвращает его в слайс
	first, err := strconv.ParseFloat(seq[pos-1], 64)
	if err != nil {
		return "", fmt.Errorf("error error in conversion %v", err)
	}
	second, err1 := strconv.ParseFloat(seq[pos+1], 64)
	if err1 != nil {
		return "", fmt.Errorf("error error in conversion %v", err)
	}
	var result float64

	switch sign {
	case "*":
		result = first * second
	case "/":
		if second == 0 {
			return "", fmt.Errorf("Devision by zero")
		}
		result = first / second
	case "-":
		result = first - second
	case "+":
		result = first + second
	}

	// составление выражения с включением результата выражения выражением
	out := strconv.FormatFloat(result, 'f', -1, 64)

	return out, nil
}

func SearchingForExpByPriority(seq []string) (string, error) {
	// поиск операций по приоритету
	for len(seq) != 1 {
		// Выполнение приоритетных операций
		for i := 0; i < len(seq); i++ {
			if string(seq[i]) == "*" || string(seq[i]) == "/" {
				resSimpleSeq, err := ExecuteBinOps(seq, i, string(seq[i])) // seq, индекс операции, операция
				if err != nil {
					return "", fmt.Errorf(err.Error())
				}
				var tempSeq = []string{}
				tempSeq = append(tempSeq, seq[:i-1]...)
				tempSeq = append(tempSeq, resSimpleSeq)
				tempSeq = append(tempSeq, seq[i+2:]...)
				i-- // уменьшаем индекс, так как длина seq изменилась
				seq = tempSeq
			}
		}

		// Выполнение менее приоритетных операций
		for i := 0; i < len(seq); i++ {
			if string(seq[i]) == "+" || string(seq[i]) == "-" {
				resSimpleSeq, err := ExecuteBinOps(seq, i, string(seq[i]))
				if err != nil {
					return "", fmt.Errorf(err.Error())
				}

				var tempSeq []string
				tempSeq = append(tempSeq, seq[:i-1]...)
				tempSeq = append(tempSeq, resSimpleSeq)
				tempSeq = append(tempSeq, seq[i+2:]...)
				i-- // уменьшаем индекс, так как длина seq изменилась
				seq = tempSeq
			}
		}
	}
	return seq[0], nil
}

func IsExpContainBrackets(exp []string) bool {
	// Проверка на содержание скобок
	for _, val := range exp {
		if val == ")" || val == "(" {
			return true
		}
	}
	return false
}

func SolveExpression(exp []string) (float64, error) {
	// основная функция решения всего выражения
	for len(exp) != 1 {
		if IsExpContainBrackets(exp) {
			indexLeftBracket := -1
			indexRightBracket := -1

			for i, val := range exp {
				if val == "(" {
					indexLeftBracket = i
				} else if val == ")" && indexLeftBracket != -1 {
					indexRightBracket = i
					break
				}
			}
			if indexLeftBracket == -1 || indexRightBracket == -1 {
				return 0.0, fmt.Errorf("incorrect sequence of parentheses")
			}

			tempExp := exp[indexLeftBracket+1 : indexRightBracket] // передача выражения вместе со скобками
			resultExp, err := SearchingForExpByPriority(tempExp)
			if err != nil {
				return 0.0, fmt.Errorf(err.Error())
			}
			var tempExpression []string
			tempExpression = append(tempExpression, exp[:indexLeftBracket]...)
			tempExpression = append(tempExpression, resultExp)
			tempExpression = append(tempExpression, exp[indexRightBracket+1:]...)

			exp = tempExpression
		} else {
			break

		}
	}
	tempExp, err := SearchingForExpByPriority(exp)
	if err != nil {
		return 0.0, fmt.Errorf("Error in conversation: %v", err)
	}
	result, _ := strconv.ParseFloat(tempExp, 64)
	return result, nil
}

func StrToSlice(str string) ([]string, error) {
	// преобразование строки в слайс
	result := []string{}
	tempNum := []string{}
	lenghtSeq := len(str)

	for i, value := range str {
		if strings.Contains("+)(-/*", string(value)) {
			if len(tempNum) > 0 {
				num := strings.Join(tempNum, "")
				result = append(result, num)
				tempNum = []string{}
			}
			result = append(result, string(value))
		} else if strings.Contains("1234567890.", string(value)) {
			tempNum = append(tempNum, string(value))

			if i == lenghtSeq-1 {
				num := strings.Join(tempNum, "")
				result = append(result, num)
			}
		} else if string(value) == " " {
			if len(tempNum) > 0 {
				num := strings.Join(tempNum, "")
				result = append(result, num)
				tempNum = []string{}
			}
		} else {
			return []string{}, fmt.Errorf("Invalid sign")
		}
	}

	return result, nil
}

func IsRightSequence(seq []string) (bool, error) {
	prevSign := string(seq[0])
	length_seq := len(seq)

	for i := 1; i < len(seq); i++ {
		if strings.Contains("*/+-", prevSign) && strings.Contains("*/+-", string(seq[i])) {
			return false, fmt.Errorf("There are two operators in a row")
		}
		if strings.Contains("1234567890.", prevSign) && strings.Contains("1234567890.", string(seq[i])) {
			return false, fmt.Errorf("There are two operands in a row")
		}
		prevSign = string(seq[i])
	}
	if strings.Contains("*/+-", string(seq[0])) {
		return false, fmt.Errorf("The expression begins with an operation")
	}
	if strings.Contains("*/+-", string(seq[length_seq-1])) {
		return false, fmt.Errorf("The expression ends with an operation")
	}
	return true, nil
}

func Calc(expression string) (float64, error) {

	if strings.Contains(expression, ")") {
		if strings.Count(expression, ")") != strings.Count(expression, ")") {
			return 0.0, fmt.Errorf("Different number of brackets")
		}
	}
	parts, err := StrToSlice(expression)
	if err != nil {
		return 0.0, fmt.Errorf(err.Error())
	}

	if len(parts) < 3 {
		return 0.0, fmt.Errorf("Incorrect Sequence")
	}
	_, err = IsRightSequence(parts)
	if err != nil {
		return 0.0, fmt.Errorf(err.Error())
	}

	result, err1 := SolveExpression(parts)
	if err1 != nil {
		return 0.0, fmt.Errorf(err1.Error())
	}
	return result, nil
	return 0.0, fmt.Errorf("Sequence is empty")
}

func main() {
	// exp := "3 + 1"
	// exp := "(10 * 3) + 5"
	// exp := "2 / 5"
	exp := "3 * 6"
	fmt.Println(Calc(exp))
}
