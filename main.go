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
	if len(os.Args) != 3 {
		fmt.Println("Использование: go run main <config.json> <events_file>")
		return
	}

	configFile := os.Args[1]
	eventsFile := os.Args[2]

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	notifier := &observer.Notifier{}

	log := logger.NewLogger()
	hand := handler.NewTracker(conf)
	gen := generator.NewEventGenerator(conf, notifier)

	notifier.Register(hand)
	notifier.Register(log)
	notifier.Register(gen)

	processEvents(eventsFile, notifier)

	hand.FinalReport()
}

func processEvents(filename string, notifier *observer.Notifier) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("ошибка открытия файла: %s", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("ошибка закрытия файла: %s", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		event, err := utils.ParseStrToEventData(line)
		if err != nil {
			fmt.Printf("Ошибка в строке '%s': %v\n", line, err)
			continue
		}
		notifier.NotifyAll(event)
	}
}
