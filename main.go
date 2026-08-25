package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"task_tracker/pkg/model"
	"task_tracker/pkg/service"
	"task_tracker/pkg/store"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args
	args = slices.Delete(args, 0, 1)

	if len(args) < 1 {
		log.Fatal("Invalid input")
	}

	command := model.NewCommand(args[0], args[1:])

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mode := os.Getenv("MODE")

	var handler store.DataHandler
	switch mode {
	case "file":
		handler, err = store.NewFileHandler()
	case "db":
		connStr := os.Getenv("DB_CONNECT")
		handler, err = store.NewDBHandler(connStr)
	default:
		fmt.Println("Invalid mode")
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}

	cr := service.NewCreator(handler)

	job, err := cr.NewJob(command)
	if err != nil {
		if errors.Is(err, service.InvalidJobError) { // test error wrapping
			log.Fatalf("Validation error: %s", err)
		}
		log.Fatal(err)
	}

	output, err := job.Execute(command.Params())
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, output)
}
