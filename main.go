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

	"github.com/kortschak/zalgo"
	"github.com/olekukonko/tablewriter"
)

// ErrNoPipe represents the error that occurs when a user doesn't use a pipe.
var ErrNoPipe = errors.New("No pipe was supplied")

// Flags.
var (
	days            string
	plain           bool
	z               bool
	delim           string
	split           string
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

	flag.StringVar(&days, "days", "weekdays", "what mapping to use, either all, weekdays or weekends")
	flag.IntVar(&timetableAmount, "amount", 0, "how many days to create the timetable for")
	flag.IntVar(&timetableEach, "each", 3, "how many items per day")
	flag.IntVar(&timetableStart, "start", 0, "the day to start on (0 indexed)")

	flag.StringVar(&delim, "delim", ", ", "what delimeter to use")
	flag.StringVar(&split, "split", "\n", "what to split the input on")
	flag.BoolVar(&plain, "plain", false, "make the output plain (no fancy table)")
	flag.BoolVar(&z, "zalgo", false, "he comes")

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
	case "weekdays", "weekday":
		dayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	case "weekends", "weekend":
		dayNames = []string{"Saturday", "Sunday"}
	case "all":
		dayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	default:
		fmt.Printf("Unknown days option %q. Please use 'weekdays' or 'all'.\n", days)
		os.Exit(1)
	}

	if timetableAmount == 0 {
		timetableAmount = len(dayNames)
	}

	rawData := timetable(
		strings.Split(subjects, split),
		timetableAmount,
		timetableEach,
	)

	output := &strings.Builder{}

	if !plain {

		table := tablewriter.NewWriter(output)
		table.SetHeader([]string{"Day", "Items"})

		var index int
		for i, val := range rawData {
			index = (i + timetableStart) % len(dayNames)
			table.Append([]string{dayNames[index], strings.Join(val, delim)})
		}

		table.SetBorder(false)
		table.SetColWidth(100)

		table.Render()

	} else {

		var index int
		for i, val := range rawData {
			index = (i + timetableStart) % len(dayNames)
			output.WriteString(fmt.Sprintf("%v: %v\n", dayNames[index], strings.Join(val, delim)))
		}

	}

	if z {
		c := zalgo.NewCorrupter(os.Stdout)
		c.Zalgo = func(n int, r rune, z *zalgo.Corrupter) bool {
			z.Up += 0.1
			z.Middle += complex(0.01, 0.01)
			z.Down += complex(real(z.Down)*0.1, 0)
			return false
		}

		c.Up = complex(0, 0.2)
		c.Middle = complex(0, 0.2)
		c.Down = complex(0.001, 0.3)

		fmt.Fprintln(c, output.String())
	} else {

		fmt.Print(output.String())
	}
}
