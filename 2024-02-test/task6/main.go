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

func startTask(r *bufio.Reader, w *bufio.Writer) error {
	var err error

	// Читаем количество наборов данных
	cntStr, err := readCountDataSet(r, 0, 10000)
	if err != nil {
		return err
	}

	for i := 0; i < cntStr; i++ {
		data, err := readLine(r, 1, 100)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		trm := newTerminal()
		for j, _ := range data {
			if err = trm.Input(data[j : j+1]); err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(w, trm.Print())
		_, err = fmt.Fprintln(w, "-")
	}

	return nil
}

/*
Объекты, используемые для решения задачи
*/

// terminal - объект для управления вводом/выводом
type terminal struct {
	screen   []string
	preview  string
	cursor   *screenCursor
	commands map[string]terminalCommand
}

type terminalCommand = func() error

func newTerminal() *terminal {
	t := &terminal{
		screen: []string{""},
		cursor: &screenCursor{
			symbol: "_",
		},
		commands: make(map[string]terminalCommand),
	}

	t.commands["L"] = t.left
	t.commands["R"] = t.right
	t.commands["U"] = t.up
	t.commands["D"] = t.down
	t.commands["B"] = t.home
	t.commands["E"] = t.end
	t.commands["N"] = t.enter

	return t
}

func (t *terminal) left() error {
	if t.cursor.y == 0 {
		return nil
	}
	t.cursor.y--

	return nil
}

func (t *terminal) right() error {
	l := t.getCurrentLineLength()

	if t.cursor.y >= l {
		t.cursor.y = l
		return nil
	}
	t.cursor.y++

	return nil
}

func (t *terminal) up() error {
	if t.cursor.x == 0 {
		return nil
	}
	t.cursor.x--

	// Если курсор выходит за пределы новой строки, то перемещаем его в конец строки
	l := t.getCurrentLineLength()
	if l < t.cursor.y {
		t.cursor.y = l
	}

	return nil
}

func (t *terminal) down() error {
	if t.cursor.x == len(t.screen)-1 {
		return nil
	}
	t.cursor.x++

	// Если курсор выходит за пределы новой строки, то перемещаем его в конец строки
	l := t.getCurrentLineLength()
	if l < t.cursor.y {
		t.cursor.y = l
	}

	return nil
}

func (t *terminal) home() error {
	t.cursor.y = 0

	return nil
}

func (t *terminal) end() error {
	t.cursor.y = t.getCurrentLineLength()

	return nil
}

func (t *terminal) enter() error {

	l := t.getCurrentLine()

	// Определяем какой кусок строки оставить в текущей строке, а какой перенести на новую
	var curLine, newLine string
	curLine = *l

	if t.cursor.y == t.getCurrentLineLength() {
		// Курсор в не в конце строки
		curLine = *l
	} else {
		newLine = curLine[t.cursor.y:]
		curLine = curLine[:t.cursor.y]
	}

	*l = curLine

	t.addLine(newLine)
	t.cursor.y = 0

	return nil
}

func (t *terminal) addLine(v string) *string {
	t.cursor.x++
	if t.cursor.x == len(t.screen) {
		t.screen = append(t.screen, v)
		return &t.screen[t.cursor.x]
	}

	newScreen := make([]string, 0)
	for i, s := range t.screen {
		if i == t.cursor.x {
			newScreen = append(newScreen, v)
		}
		newScreen = append(newScreen, s)
	}

	t.screen = newScreen

	return &t.screen[t.cursor.x]
}

func (t *terminal) getCurrentLine() *string {
	if t.cursor.x > len(t.screen) {
		return nil
	}
	return &t.screen[t.cursor.x]
}

func (t *terminal) getCurrentLineLength() int {
	if t.cursor.x > len(t.screen) {
		return 0
	}
	return len(t.screen[t.cursor.x])
}

func (t *terminal) inputToCurrentLine(v string) error {
	l := t.getCurrentLine()

	if t.cursor.y == t.getCurrentLineLength() {
		*l += v
		t.cursor.y++
		return nil
	}

	curVal := *l
	afterLine := curVal[t.cursor.y:]
	beforeLine := curVal[:t.cursor.y]
	*l = beforeLine + v + afterLine
	t.cursor.y++

	return nil
}

func (t *terminal) Input(v string) error {
	if cmd, ok := t.commands[v]; ok {
		err := cmd()
		t.preview = t.Screen()
		return err
	}

	if ok, err := regexp.MatchString("^[a-z,0-9]$", v); !ok || err != nil {
		if err != nil {
			return err
		}
		return errors.New("incorrect terminal input")
	}

	err := t.inputToCurrentLine(v)
	t.preview = t.Screen()
	return err
}

func (t *terminal) Print() string {
	return strings.Join(t.screen, "\n")
}

func (t *terminal) Screen() string {
	var res string
	for i, _ := range t.screen {
		if i == t.cursor.x {
			curVal := t.screen[i]
			res += curVal[:t.cursor.y] + t.cursor.symbol + curVal[t.cursor.y:] + "\n"
		} else {
			res += t.screen[i] + "\n"
		}
	}
	return res
}

// screenCursor - объект для обозначения текущего положения курсора
type screenCursor struct {
	symbol string
	x      int
	y      int
}

/*
Считывание данных
*/

// readCountDataSet - считывание количества наборов данных
func readCountDataSet(r *bufio.Reader, min int, max int) (int, error) {
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

// readLine - чтение строки с данными
func readLine(r *bufio.Reader, min int, max int) (string, error) {
	lBytes, _, err := r.ReadLine()
	l := string(lBytes)
	if Debug {
		println("line", l)
	}
	if err != nil {
		return "", fmt.Errorf("incorrect input: %v", err.Error())
	}
	l = strings.TrimSpace(l)
	if len(l) < min || len(l) > max {
		return "", errors.New("incorrect data")
	}

	return l, nil
}
