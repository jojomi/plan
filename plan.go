package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Priority uint8

const (
	VeryHigh Priority = iota
	High
	Standard
	Low
	VeryLow
)

type Status uint8

const (
	Open Status = iota
	InProgress
	Done
)

// Task is a task.
type Task struct {
	ID       string
	Title    string
	Text     string
	Duration float64
	Priority Priority
	Status   Status
}

var pause float64 = 10

func main() {
	tasks := []Task{}

	// Parse
	file, err := os.Open("input.csv")
	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("Error:", err)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#L94
	//reader.Comma = ','
	lineCount := 0
	for {
		lineCount++
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
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
			task.Priority = VeryHigh
		case "high":
			task.Priority = High
		case "standard":
			task.Priority = Standard
		case "low":
			task.Priority = Low
		case "very low":
			task.Priority = VeryLow
		}

		switch strings.ToLower(record[5]) {
		case "open":
			task.Status = Open
		case "inprogress":
			task.Status = InProgress
		case "done":
			task.Status = Done
		}
		tasks = append(tasks, task)
	}
	// Sort

	// Output
	now := time.Now()
	start := now
	const layout = "15:04"
	for _, task := range tasks {
		if task.Status == Done {
			continue
		}
		start = start.Add(time.Duration(task.Duration) * time.Minute).Add(time.Duration(pause) * time.Minute)
		marker := ""
		if task.Priority >= High {
			marker = "! "
		}
		fmt.Printf("%s%s – %s [%.0f minutes, due at %s]\n", marker, task.Title, task.Text, task.Duration, start.Format(layout))
	}
}
