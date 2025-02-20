package config

import (
	"os"
	"reflect"
	"strconv"
)

var (
	TIME_ADDITION_MS        = 10
	TIME_SUBTRACTION_MS     = 10
	TIME_MULTIPLICATIONS_MS = 10
	TIME_DIVISIONS_MS       = 10
	COMPUTING_POWER         = 10
)

type Config struct {
	Addr                    string
	TIME_ADDITION_MS        int
	TIME_SUBTRACTION_MS     int
	TIME_MULTIPLICATIONS_MS int
	TIME_DIVISIONS_MS       int
	COMPUTING_POWER         int
}

func ConfigFromEnv() *Config {
	Config := new(Config)

	Config.Addr = os.Getenv("PORT")
	if Config.Addr == "" {
		Config.Addr = "8080"
	}

	// Создаем мапу для соответствия имен полей и значений
	defaultValues := map[string]int{
		"TIME_ADDITION_MS":        TIME_ADDITION_MS,
		"TIME_SUBTRACTION_MS":     TIME_SUBTRACTION_MS,
		"TIME_MULTIPLICATIONS_MS": TIME_MULTIPLICATIONS_MS,
		"TIME_DIVISIONS_MS":       TIME_DIVISIONS_MS,
		"COMPUTING_POWER":         COMPUTING_POWER,
	}

	// Используем reflect для итерации по полям структуры
	val := reflect.ValueOf(Config).Elem()
	typ := reflect.TypeOf(*Config)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Проверяем, является ли поле целым числом и равно ли оно нулю
		if field.Kind() == reflect.Int && field.Int() == 0 {
			// Получаем значение по имени поля из мапы
			if defaultValue, exists := defaultValues[fieldType.Name]; exists {
				field.SetInt(int64(defaultValue))
				os.Setenv(fieldType.Name, strconv.Itoa(defaultValue))
			}
		}
	}

	return Config
}
