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

func startTask(r *bufio.Reader, w *bufio.Writer) error {
	var err error

	// Читаем количество наборов данных
	cntStr, err := readCountDataSet(r, 0, 100)
	if err != nil {
		return err
	}

	for i := 0; i < cntStr; i++ {
		cntPages, err := readCountPages(r, 2, 100)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		data, err := readLine(r, 1, 100)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		prn := newPrintTask(cntPages)
		for j, _ := range data {
			err = prn.AddPrintedPages(data[j])
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(w, prn.Print())
	}

	return nil
}

/*
Объекты, используемые для решения задачи
*/

// printTask - объект для формирования задания на печать
type printTask struct {
	printedPages map[int]bool
	pageCount    int
}

func newPrintTask(cnt int) *printTask {
	t := &printTask{
		printedPages: make(map[int]bool),
		pageCount:    cnt,
	}
	for i := 1; i <= cnt; i++ {
		t.printedPages[i] = false
	}
	return t
}

func (t *printTask) AddPrintedPages(s string) error {
	p := newPageInterval(0)
	err := p.Parse(s)
	if err != nil {
		return err
	}

	for i := p.first; i <= p.last; i++ {
		t.printedPages[i] = true
	}

	return nil
}

func (t *printTask) Print() string {
	printInterval := make([]*pageInterval, 0)
	var curInterval *pageInterval
	for i := 1; i <= t.pageCount; i++ {
		if t.printedPages[i] {
			continue
		}

		if curInterval == nil {
			curInterval = newPageInterval(i)
			printInterval = append(printInterval, curInterval)
			continue
		}

		if curInterval.Expand(i) {
			continue
		}

		curInterval = newPageInterval(i)
		printInterval = append(printInterval, curInterval)
	}

	pr := make([]string, 0)
	for _, v := range printInterval {
		pr = append(pr, v.ToString())
	}

	return strings.Join(pr, ",")
}

type pageInterval struct {
	first int
	last  int
}

func newPageInterval(start int) *pageInterval {
	return &pageInterval{
		first: start,
		last:  start,
	}
}

func (t *pageInterval) Parse(s string) (err error) {
	t.first, err = strconv.Atoi(s)
	if err == nil {
		t.last = t.first
		return nil
	}

	interval := strings.Split(s, "-")
	if len(interval) != 2 {
		return errors.New("incorrect page interval")
	}

	t.first, err = strconv.Atoi(interval[0])
	if err != nil {
		return errors.New("incorrect page interval")
	}
	t.last, err = strconv.Atoi(interval[1])
	if err != nil {
		return errors.New("incorrect page interval")
	}

	return
}

func (t *pageInterval) Expand(page int) bool {
	if page-t.last != 1 {
		return false
	}

	t.last = page
	return true
}

func (t *pageInterval) ToString() string {
	if t.first == t.last {
		return strconv.Itoa(t.first)
	}

	return fmt.Sprintf("%v-%v",
		strconv.Itoa(t.first),
		strconv.Itoa(t.last),
	)
}

/*
Считывание данных
*/

// readCountDataSet - считывание количества наборов данных
func readCountDataSet(r *bufio.Reader, min int, max int) (int, error) {
	cntStr, _ := r.ReadString('\n')
	if Debug {
		println("count data set", cntStr)
	}
	cntStr = strings.TrimSpace(cntStr)
	cnt, err := strconv.Atoi(cntStr)
	if err != nil || cnt < min || cnt > max {
		return 0, errors.New("incorrect count line")
	}
	return cnt, nil
}

// readCountPages - считывание количества страниц к печати
func readCountPages(r *bufio.Reader, min int, max int) (int, error) {
	cntStr, _ := r.ReadString('\n')
	if Debug {
		println("count pages", cntStr)
	}
	cntStr = strings.TrimSpace(cntStr)
	cnt, err := strconv.Atoi(cntStr)
	if err != nil || cnt < min || cnt > max {
		return 0, errors.New("incorrect count line")
	}
	return cnt, nil
}

// readLine - чтение строки с данными
func readLine(r *bufio.Reader, min int, max int) ([]string, error) {
	lBytes, _, err := r.ReadLine()
	l := string(lBytes)
	if Debug {
		println("line", l)
	}
	if err != nil {
		return nil, fmt.Errorf("incorrect input: %v", err.Error())
	}
	l = strings.TrimSpace(l)

	dataStr := strings.Split(l, ",")

	if len(dataStr) < min || len(dataStr) > max {
		return nil, errors.New("incorrect data")
	}

	return dataStr, nil
}
