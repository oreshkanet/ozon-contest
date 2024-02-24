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

type conditioner struct {
	minT int
	maxT int
}

func newConditioner(minT int, maxT int) *conditioner {
	return &conditioner{
		minT: minT,
		maxT: maxT,
	}
}

func (c *conditioner) correctTemp(cmd string, temp string) error {

	t, err := strconv.Atoi(temp)
	if err != nil {
		return err
	}

	switch cmd {
	case ">=":
		if t > c.minT {
			c.minT = t
		}
	case "<=":
		if t < c.maxT {
			c.maxT = t
		}
	default:
		return errors.New("incorrecr command")
	}

	return nil
}

func (c *conditioner) isPossible() bool {
	return c.minT <= c.maxT
}

func (c *conditioner) printIsPossibleString() string {
	if !c.isPossible() {
		return "-1"
	}
	return fmt.Sprint(c.minT)
}

func startTask(r *bufio.Reader, w *bufio.Writer) error {
	var err error
	cntStr, err := readCountLines(r, 0, 10000)
	if err != nil {
		return err
	}

	for i := 0; i < cntStr; i++ {
		cntEmp, err := readCountEmployees(r, 0, 10000)
		if err != nil {
			if Debug {
				println("err line employees", err.Error())
			}
			return err
		}

		cnd := newConditioner(15, 30)
		for j := 0; j < cntEmp; j++ {
			cmd, err := readLine(r)
			if err != nil {
				if Debug {
					println("err line", err.Error())
				}
				return err
			}

			err = cnd.correctTemp(cmd[0], cmd[1])
			if err != nil {
				if Debug {
					println("err command", err.Error())
				}
				return err
			}
			_, err = fmt.Fprintln(w, cnd.printIsPossibleString())
		}
		_, err = fmt.Fprintln(w, "")
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

func readCountEmployees(r *bufio.Reader, min int, max int) (int, error) {
	cntStr, _ := r.ReadString('\n')
	if Debug {
		println("count employees", cntStr)
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
	cmd := strings.Split(l, " ")
	if len(cmd) != 2 {
		return nil, fmt.Errorf("incorrect input")
	}
	return cmd, nil
}
