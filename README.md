# Timetable
Timetable is a simple tool for generating a timetable. My use case is to generate a revision timetable.

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
     DAY    |         ITEMS           
------------+-------------------------
  Monday    | Maths, Art, Drama       
  Tuesday   | English, IT, History    
  Wednesday | Science, Geography, PE  
  Thursday  | Music, Art, Maths       
  Friday    | English, IT, History    
```

## "Advanced" Usage
Because why not ruthlessly overengineer what should be a super simple program?

* `-amount`: **int**

  how many days to create the timetable for
  
* `-days`: **string**

  what mapping to use, either all or weekdays (default "weekdays")
  
* `-delim`: **string**

  what delimeter to use (default ", ")
  
* `-each`: **int**

  how many items per day (default 3)
  
* `-plain`:

  make the output plain
  
* `-split`: **string**

  what to split the input on (default "\n")
  
* `-start`: **int**

  the day to start on (0 indexed)
  