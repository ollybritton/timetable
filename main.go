package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ErrNoPipe represents the error that occurs when a user doesn't use a pipe.
var ErrNoPipe = errors.New("No pipe was supplied")

// Flags.
var (
	days            string
	timetableAmount int
	timetableEach   int
	timetableStart  int
)

func getPipe() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return "", ErrNoPipe
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return string(output), nil
}

func shuffle(list []string) []string {
	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i)

		temp := list[j]
		list[j] = list[i]
		list[i] = temp
	}

	return list
}

func slightlyShuffle(list []string) []string {
	newlist := make([]string, len(list))
	copy(newlist, list)

	// swaps list items from left to right 3/4ths of the time
	for i := 0; i < len(newlist)-1; i++ {
		if rand.Intn(3) > 1 {
			temp := newlist[i]
			newlist[i] = newlist[i+1]
			newlist[i+1] = temp
		}
	}

	return newlist
}

func chunk(list []string, n int) [][]string {
	chunks := [][]string{}

	for i := 0; i < len(list); i += n {
		if i+n > len(list) {
			return append(chunks, list[i:])
		}

		chunks = append(chunks, list[i:i+n])
	}

	return chunks
}

func timetable(list []string, length int, each int) [][]string {
	table := [][]string{}
	full := []string{}

	list = shuffle(list)

	// In order to have l chunks of e size, we need a list l*e long.
	for len(full) < length*each {
		list = slightlyShuffle(list)
		full = append(full, list...)
	}

	chunks := chunk(full, each)
	for _, c := range chunks {
		table = append(table, c)
	}

	return table[:length]
}

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&days, "days", "weekdays", "what mapping to use, either all or weekdays")
	flag.IntVar(&timetableAmount, "amount", 5, "how many days to create the timetable for")
	flag.IntVar(&timetableEach, "each", 3, "how many items per day")
	flag.IntVar(&timetableStart, "start", 0, "the day to start on (0 indexed)")

	flag.Parse()
}

func main() {
	subjects, err := getPipe()
	if errors.Is(err, ErrNoPipe) {
		fmt.Println("This command is designed to work with pipes.")
		fmt.Println("Example: cat subjects.txt | timetable")
		os.Exit(1)
	}

	if timetableAmount < 0 || timetableEach < 0 || timetableStart < 0 {
		fmt.Println("Invalid negative argument")
		os.Exit(1)
	}

	var dayNames []string
	switch days {
	case "weekdays":
		dayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	case "all":
		dayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	default:
		fmt.Printf("Unknown days option %q. Please use 'weekdays' or 'all'.\n", days)
		os.Exit(1)
	}

	table := timetable(
		strings.Split(subjects, "\n"),
		timetableAmount,
		timetableEach,
	)

	for i, val := range table {
		fmt.Println(dayNames[(i+timetableStart)%len(dayNames)], "=>", strings.Join(val, ", "))
	}
}
