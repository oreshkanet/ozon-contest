package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
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

type carNumber struct {
	length int
	regexp *regexp.Regexp
}

func newCarNumber(length int, expr string) *carNumber {
	r, err := regexp.Compile(expr)
	if err != nil {
		log.Fatal(err)
	}
	return &carNumber{
		length: length,
		regexp: r,
	}
}

func (c *carNumber) Parse(s string) (string, string, bool) {
	if len(s) < c.length {
		return s, "", false
	}
	n := s[:c.length]

	if !c.regexp.MatchString(n) {
		return s, "", false
	}

	s = s[c.length:]

	return s, n, true
}

func startTask(r *bufio.Reader, w *bufio.Writer) error {

	cntStr, err := readCountLines(r, 0, 1000)
	if err != nil {
		return err
	}

	templates := []*carNumber{
		newCarNumber(4, `^[A-Z][0-9][A-Z]{2}$`),
		newCarNumber(5, `^[A-Z][0-9]{2}[A-Z]{2}$`),
	}

	for i := 0; i < cntStr; i++ {
		carNumbers := make([]string, 0)
		rawNumbers, err := readLine(r)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		var ok bool
		carNumbers, ok, err = parseNumbers(carNumbers, rawNumbers, templates)
		if err != nil {
			if Debug {
				println("err parse", err.Error())
			}
			return err
		}
		if !ok {
			_, err = fmt.Fprintln(w, "-")
		} else {
			_, err = fmt.Fprintln(w, strings.Join(carNumbers, " "))
		}
		if err != nil {
			if Debug {
				println("err print", err.Error())
			}
			return err
		}
	}

	return nil
}

func parseNumbers(result []string, s string, templates []*carNumber) ([]string, bool, error) {
	var err error
	for _, tmpl := range templates {
		var num string
		var ok bool
		s, num, ok = tmpl.Parse(s)
		if !ok {
			continue
		}
		result = append(result, num)
		if s == "" {
			// Успешно дошли до конца строки
			return result, true, nil
		}

		// Парсим следующий номер
		result, ok, err = parseNumbers(result, s, templates)
		if err != nil {
			return result, false, err
		}
		if ok {
			// Все последующие номера успешно распарсили
			return result, true, nil
		}
		// Проверяем следующий шаблон номеров
	}
	return result, false, nil
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

func readLine(r *bufio.Reader) (string, error) {
	lBytes, _, err := r.ReadLine()
	l := string(lBytes)
	if Debug {
		println("line", l)
	}
	if err != nil {
		return "", fmt.Errorf("incorrect input: %v", err.Error())
	}
	l = strings.TrimSpace(l)
	if len(l) == 0 || len(l) > 50 {
		return "", fmt.Errorf("incorrect input")
	}
	return l, nil
}
