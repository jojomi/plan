package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type priority uint8

const (
	veryHigh priority = iota
	high
	standard
	low
	veryLow
)

type status uint8

const (
	open status = iota
	inProgress
	done
)

// Task is a task.
type Task struct {
	ID       string
	Title    string
	Text     string
	Duration float64
	Priority priority
	Status   status
}

var pause time.Duration = 5 * time.Minute

func main() {
	tasks := []Task{}

	// Parse
	filename := "input.csv"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineCount := 0
	for {
		lineCount++
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if lineCount == 1 {
			continue
		}

		task := Task{}
		task.ID = record[0]
		task.Title = record[1]
		task.Text = record[2]
		dur, _ := strconv.ParseFloat(record[3], 64)
		task.Duration = dur

		switch strings.ToLower(record[4]) {
		case "very high":
			task.Priority = veryHigh
		case "high":
			task.Priority = high
		case "standard":
			task.Priority = standard
		case "low":
			task.Priority = low
		case "very low":
			task.Priority = veryLow
		}

		switch strings.ToLower(record[5]) {
		case "open":
			task.Status = open
		case "inprogress":
			task.Status = inProgress
		case "done":
			task.Status = done
		}
		tasks = append(tasks, task)
	}
	// TODO sort

	// output
	now := time.Now()
	start := now
	const layout = "15:04"
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	if len(tasks) > 0 {
		fmt.Fprintln(w, "Prio\tTitle\tText\tDuration\tEnd\t")
	}
	for _, task := range tasks {
		if task.Status == done {
			continue
		}
		start = start.Add(time.Duration(task.Duration) * time.Minute).Add(pause)
		marker := ""
		if task.Priority >= high {
			marker = "! "
		}
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%.0f Minutes\t%s\t", marker, task.Title, task.Text, task.Duration, start.Format(layout)))
	}
	w.Flush()
}
