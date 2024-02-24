package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var Debug bool = false

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	err := startTask(r, w)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err = w.Flush(); err != nil {
		log.Fatal(err.Error())
	}
}

type dateTime struct {
	day   string
	month string
	year  string
}

func newDateTime(day string, month string, year string) *dateTime {
	return &dateTime{
		day:   day,
		month: month,
		year:  year,
	}
}

func (d *dateTime) isCorrect() bool {
	var err error
	var day, month, year int
	if day, err = strconv.Atoi(d.day); err != nil {
		return false
	}
	if month, err = strconv.Atoi(d.month); err != nil {
		return false
	}
	if year, err = strconv.Atoi(d.year); err != nil {
		return false
	}
	_, err = time.Parse(
		time.DateOnly,
		fmt.Sprintf("%d-%02d-%02d", year, month, day),
	)
	if err != nil {
		return false
	}

	return true
}

func (d *dateTime) printCorrectString() string {
	if d.isCorrect() {
		return "YES"
	}
	return "NO"
}

func startTask(r *bufio.Reader, w *bufio.Writer) error {

	cntStr, err := readCountLines(r, 0, 1000)
	if err != nil {
		return err
	}

	for i := 0; i < cntStr; i++ {
		dt, err := readLine(r)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		date := newDateTime(dt[0], dt[1], dt[2])
		_, err = fmt.Fprintln(w, date.printCorrectString())
	}

	return nil
}

func readCountLines(r *bufio.Reader, min int, max int) (int, error) {
	cntStr, _ := r.ReadString('\n')
	if Debug {
		println("count strings", cntStr)
	}
	cntStr = strings.TrimSpace(cntStr)
	cnt, err := strconv.Atoi(cntStr)
	if err != nil || cnt < min || cnt > max {
		return 0, errors.New("incorrect count line")
	}
	return cnt, nil
}

func readLine(r *bufio.Reader) ([]string, error) {
	lBytes, _, err := r.ReadLine()
	l := string(lBytes)
	if Debug {
		println("line", l)
	}
	if err != nil {
		return nil, fmt.Errorf("incorrect input: %v", err.Error())
	}
	l = strings.TrimSpace(l)
	res := strings.Split(l, " ")
	if len(res) != 3 {
		return nil, fmt.Errorf("incorrect input")
	}
	return res, nil
}
