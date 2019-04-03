package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"time"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"strconv"
)


/////////////////////////////////////////////////////////////////////////////
func main() {
	// Parse args
	//if len(os.Args) != 2 {
	//	die(1, "Usage")
	//}
	//inputFileName := os.Args[1]
	//
	//r, err := os.Open(inputFileName)
	//if err != nil {
	//	die(1,"ERROR! Can't open input file")
	//}

	// Get I/O streams
	lineScanner := bufio.NewScanner(os.Stdin) // TODO bufio.NewScanner(getInputReader())
	lineWriter := os.Stdout
	progressWriter := os.Stderr // TODO if -v

	// Get the throttle
	throttle := getThrottle()

	// Main loop
	lineNo := 0
	progressPrevLen := 0
	for lineScanner.Scan(){
		// Read line
		lineNo++
		line := lineScanner.Text()

		// Throttle and print out progress
		wait := true
		msg := ""
		for wait {
			wait, msg = throttle.poll()

			// Print progress
			progressText := fmt.Sprintf("[%05d][%s] %s", lineNo, msg, line)
			if len(progressText) > 80 {
				progressText = progressText[:79]
			}
			fmt.Fprint(progressWriter, "\r" + strings.Repeat(" ", progressPrevLen))
			fmt.Fprint(progressWriter, "\r" + progressText)
			progressPrevLen = len(progressText)
		}
		fmt.Fprintln(progressWriter)

		// Write out the line
		fmt.Fprintln(lineWriter, line)
		// TODO flush
	}

	// Exit
	os.Exit(0)
}

/////////////////////////////////////////////////////////////////////////////
func die (code int, msg string){
	fmt.Println(msg)
	os.Exit(code)
}

type throttle interface {
	init()
	poll() (bool, string)
}

func getThrottle() throttle {
	//return timeThrottle{
	//	delay: 1000 * time.Millisecond,
	//}

	t := innodbThrottle{}
	t.init()
	return &t
}

/////////////////////////////////////////////////////////////////////////////
type timeThrottle struct {
	delay time.Duration
}

func (t timeThrottle) init() {
}

func (t timeThrottle) poll() (bool, string) {
	time.Sleep(t.delay)
	return false, fmt.Sprintf("%fs", t.delay.Seconds())
}


/////////////////////////////////////////////////////////////////////////////
type innodbThrottle struct {
	db *sql.DB
	status map[string]int
}

// Connect to MySQL
func (i *innodbThrottle) init() {
	i.status = make(map[string]int)

	// MySQL connection params
	conf := mysql.Config{
		User: "root",
		Passwd: "",
		//Net: "tcp",
		//Addr: "localhost",
		Net: "unix",
		Addr: "/tmp/mysql.sock",
		AllowNativePasswords: true,
	}
	dsn := conf.FormatDSN()

	// Connect to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		die(1, fmt.Sprint(err))
	}
	i.db = db

}

// Pull status variables from MySQL
func (i *innodbThrottle) update(){
	// Get status
	rows, err := i.db.Query("SHOW STATUS LIKE '%innodb%'")
	if err != nil {
		die(1, fmt.Sprintf("DB ERROR: %s", err))
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			value string
		)
		if err := rows.Scan(&name, &value); err != nil {
			die(1, fmt.Sprint(err))
		}
		//fmt.Printf("%20s %s\n", name, value)
		intval, err := strconv.Atoi(value)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "ERROR: could not parse \"%s\" value \"%s\" as int: %s\n", name, value, err)
			//i.status[name] = -1
		} else {
			i.status[name] = intval
		}
	}
}

func (i *innodbThrottle) poll() (bool, string){
	i.update()
	// TODO
	time.Sleep(500 * time.Millisecond)
	key := "Innodb_pages_read"
	return false, fmt.Sprintf("%s: %d", key, i.status[key])
}

/////////////////////////////////////////////////////////////////////////////
//func getInputReader() io.Reader {
//func getInputReader() io.Reader {
//	if len(os.Args) == 1 {
//		return os.Stdin
//	}
//
//	if len(os.Args) != 2 {
//		die(1, "Usage")
//	}
//	inputFileName := os.Args[1]
//
//	r, err := os.Open(inputFileName)
//	if err != nil {
//		die(1,"ERROR! Can't open input file")
//	}
//
//	return r
//}
