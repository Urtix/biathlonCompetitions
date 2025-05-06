package main

import (
	"biathlonCompetitions/internal/config"
	"biathlonCompetitions/internal/generator"
	"biathlonCompetitions/internal/handler"
	"biathlonCompetitions/internal/logger"
	"biathlonCompetitions/internal/observer"
	"biathlonCompetitions/internal/utils"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Парсинг аргументов командной строки
	cfgPath, eventsPath, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Загрузка конфигурации из файла
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Инициализация системы уведомлений (реализация паттерна Observer)
	notifier := &observer.Notifier{}

	// Создание основных компонентов системы
	log := logger.NewLogger()                         // Логирование событий
	hand := handler.NewHandler(cfg)                   // Обработчик событий
	gen := generator.NewEventGenerator(cfg, notifier) // Генератор новых событий

	// Регистрация компонентов как подписчиков уведомлений
	notifier.Register(hand)
	notifier.Register(log)
	notifier.Register(gen)

	// Обработка событий из файла
	processEvents(eventsPath, notifier)

	// Вывод финального отчета после обработки всех событий
	hand.FinalReport()
}

// Парсинг аргументов командной строки
// Ожидается два аргумента: путь к конфигу и путь к файлу событий
func parseArgs() (string, string, error) {
	if len(os.Args) != 3 {
		return "", "", fmt.Errorf("use: %s <config.json> <events_file>", os.Args[0])
	}
	return os.Args[1], os.Args[2], nil
}

// Обработка событий из файла
func processEvents(filename string, notifier *observer.Notifier) {
	// Открытие файла с событиями
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error open file: %s", err)
		return
	}
	// Гарантированное закрытие файла при выходе из функции
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("Error close file: %s", err)
		}
	}(file)

	// Построчное чтение файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Парсинг строки в структуру события
		event, err := utils.ParseStrToEventData(line)
		if err != nil {
			fmt.Printf("Error in line '%s': %v\n", line, err)
			continue
		}
		// Рассылка события всем подписчикам
		notifier.NotifyAll(event)
	}
}
