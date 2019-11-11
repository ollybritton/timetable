# Timetable
Timetable is a ~~simple~~ needlessly overengineered tool for generating a timetable. My use case is to generate a revision timetable.

## Installation
Assuming you have the `go` command installed:
```
go get -u github.com/ollybritton/timetable
```

## Usage
To generate a timetable from a newline-seperated list of items:

```bash
$ cat items.txt | timetable
```

Example output:

```
Seed: 1573504119091066000

     DAY    |         ITEMS           
------------+-------------------------
  Monday    | Maths, Art, Drama       
  Tuesday   | English, IT, History    
  Wednesday | Science, Geography, PE  
  Thursday  | Music, Art, Maths       
  Friday    | English, IT, History    
```

## Wtf is this about the seed
Like I said, poorly overengineered: That is the seed for the RNG. It allows you to generate timetables and then extend them later.

For example, say I make a timetable stretching from Monday to Friday and then on Friday I realise I wanted a table that lasted the full weekend. I could just generate a new timetable, but it would mean that some items in the list would be under represented. So instead, I just use the `seed` output of the first run as an input to the second:

```bash
$ cat items.txt | timetable -seed 1573504119091066000 -days all
```

Then I will continue that timetable instead of creating a new one.

Horrible, I know.

## "Advanced" Usage
Because why not ruthlessly overengineer what should be a super simple program? Run `timetable -help`

```
Usage of timetable:
  -amount int
        how many days to create the timetable for
  -days string
        what mapping to use, either all, weekdays or weekends (default "weekdays")
  -delim string
        what delimeter to use (default ", ")
  -each int
        how many items per day (default 3)
  -no-seed
        print the random number seed
  -plain
        make the output plain (no fancy table)
  -seed int
        seed to use (default 1573504447938477000)
  -split string
        what to split the input on (default "\n")
  -start int
        the day to start on (0 indexed)
  -zalgo
        he comes
```