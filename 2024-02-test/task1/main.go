package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

type fleet struct {
	ships map[string]int8
}

func newFleet() *fleet {
	return &fleet{
		ships: make(map[string]int8),
	}
}

func (f *fleet) addShip(t string) {
	f.ships[t]++
}

func (f *fleet) isCorrect() bool {
	return f.ships["1"] == 4 &&
		f.ships["2"] == 3 &&
		f.ships["3"] == 2 &&
		f.ships["4"] == 1
}

func (f *fleet) printCorrectString() string {
	if f.isCorrect() {
		return "YES"
	}
	return "NO"
}

func startTask(r *bufio.Reader, w *bufio.Writer) error {

	cntStr, err := readCountLines(r, 0, 1000)
	if err != nil {
		return err
	}

	fleets := make([]*fleet, cntStr)
	for i := 0; i < cntStr; i++ {
		fleets[i] = newFleet()
		ships, err := readLine(r)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}
		for _, ship := range ships {
			fleets[i].addShip(ship)
		}
	}

	for _, fleet := range fleets {
		_, err = fmt.Fprintln(w, fleet.printCorrectString())
		if err != nil {
			if Debug {
				println("err print", err.Error())
			}
			return err
		}
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
	return strings.Split(l, " "), nil
}
