package calculator

import (
	"calc/models"
	"fmt"
	"strconv"
	"strings"
)

// Выполняет выражение и возвращает его в слайс
func ExecuteBinOps(seq []string, pos int, sign string) (string, error) {
	first, err := strconv.ParseFloat(seq[pos-1], 64)
	if err != nil {
		return "", ErrConvertingNumberToFloatType
	}
	second, err1 := strconv.ParseFloat(seq[pos+1], 64)
	if err1 != nil {
		return "", ErrConvertingNumberToFloatType
	}
	var result float64
	// выполнений операций над двумя операндами
	switch sign {
	case "*":
		result = first * second
	case "/":
		if second == 0 {
			fmt.Println(ErrDivisionByZero)
			return "", ErrDivisionByZero
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

// Поиск операций по приоритету
func SearchingForExpByPriority(seq []string) (string, error) {

	for len(seq) != 1 {
		// Выполнение приоритетных операций
		for i := 0; i < len(seq); i++ {
			if string(seq[i]) == "*" || string(seq[i]) == "/" {
				resSimpleSeq, err := ExecuteBinOps(seq, i, string(seq[i])) // seq, индекс операции, операция

				if err != nil {
					fmt.Printf("Error from SearchingForExpByPriority: %s, %v", resSimpleSeq, err)
					return "", err
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
					return "", err
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

// Проверка на содержание скобок
func IsExpContainBrackets(exp []string) bool {
	for _, val := range exp {
		if val == ")" || val == "(" {
			return true
		}
	}
	return false
}

// основная функция решения всего выражения
func SolveExpression(exp []string) (float64, error) {
	// resultChan := make(chan Result)
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
				return 0.0, ErrIncorrectSeqOfParenthese
			}

			tempExp := exp[indexLeftBracket+1 : indexRightBracket] // передача выражения вместе со скобками
			resultExp, err := SearchingForExpByPriority(tempExp)
			if err != nil {
				return 0.0, err
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
		return 0.0, err
	}
	result, _ := strconv.ParseFloat(tempExp, 64)
	return result, nil
}

// Проверка строки на правильную последовательность выражений
func IsRightSequence(seq []string) (bool, error) {
	prevSign := string(seq[0])
	length_seq := len(seq)

	for i := 1; i < len(seq); i++ {
		if strings.Contains("*/+-(", prevSign) && strings.Contains("*/+-", string(seq[i])) {
			return false, ErrTwoOperatorsInRow
		}
		if strings.Contains("1234567890.", prevSign) && strings.Contains("1234567890.", string(seq[i])) {
			return false, ErrTwoOperandsInRow
		}
		prevSign = string(seq[i])
	}
	if strings.Contains("*/+-", string(seq[0])) {
		return false, ErrExpStartsWithOperator
	}
	if strings.Contains("*/+-", string(seq[length_seq-1])) {
		return false, ErrExpEndsWithOperator
	}
	return true, nil
}

// Преобразование строки в слайс
func StrToSlice(str string) ([]string, error) {
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
			return []string{}, ErrInvalidExpression
		}
	}
	return result, nil
}

// Основная функция расчёта
func Calc(e *models.Expression) (float64, error) {
	fmt.Println("Зашел в Calc")
	if strings.Count(e.Exp, ")") != strings.Count(e.Exp, "(") {
		return 0.0, ErrDiffNumberOfBrackets
	}
	parts, err := StrToSlice(e.Exp)
	if err != nil {
		return 0.0, err
	}

	if len(parts) < 3 {
		return 0.0, ErrInvalidExpression
	}
	_, err = IsRightSequence(parts)
	if err != nil {
		return 0.0, err
	}

	result, err1 := SolveExpression(parts)
	if err1 != nil {
		return 0.0, err
	}
	return result, nil
}

// func main() {
// 	// exp := "3 + 1"
// 	// exp := "(10 * 3) + 5"
// 	// exp := "2 / 5"
// 	exp := "3)6"
// 	fmt.Println(Calc(exp))
// }
