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
	cntStr, err := readCountDataSet(r, 0, 10000)
	if err != nil {
		return err
	}

	for i := 0; i < cntStr; i++ {
		dataLen, err := readDataLength(r, 1, 50)
		if err != nil {
			if Debug {
				println("err data length", err.Error())
			}
			return err
		}
		data, err := readLine(r, -1000, 1000)
		if err != nil {
			if Debug {
				println("err line", err.Error())
			}
			return err
		}

		cmps := newCompressor(data, dataLen)
		err = cmps.compress()
		if err != nil {
			if Debug {
				println("err compress", err.Error())
			}
			return err
		}

		_, err = fmt.Fprintln(w, strconv.Itoa(cmps.Length()))
		_, err = fmt.Fprintln(w, cmps.toString())
	}

	return nil
}

/*
Объекты, используемые для решения задачи
*/

// bit - один элемент сжатых данных
type bit struct {
	start  int
	offset int
	sign   int
}

func newBit(start int) *bit {
	return &bit{
		start: start,
	}
}

func (b *bit) GetOffset() int {
	return b.offset * b.sign
}

func (b *bit) LastValue() int {
	return b.start + b.GetOffset()
}

func (b *bit) AddIfPossible(v int) bool {
	offset := v - b.LastValue()
	if offset != -1 && offset != 1 {
		// оффсет больше 1, значит добавить нельзя
		return false
	}
	if b.sign == 0 {
		// В бите только одно значение и оффсета ещё нет - добавляем его
		b.sign = offset
	} else if b.sign != offset {
		// текущий оффсет не совпадает со знаком бита
		return false
	}

	b.offset++

	return true
}

// compressor - объект для сжатия данных
type compressor struct {
	data       []int
	dataLength int
	compressed []*bit
}

func newCompressor(data []int, len int) *compressor {
	return &compressor{
		data:       data,
		dataLength: len,
	}
}

func (c *compressor) AddBit(start int) *bit {
	b := newBit(start)
	c.compressed = append(c.compressed, b)
	return b
}

func (c *compressor) Length() int {
	return len(c.compressed) * 2
}

func (c *compressor) compress() error {
	var curBit *bit
	for _, v := range c.data {
		if curBit == nil {
			// Это первый элемент
			curBit = c.AddBit(v)
			continue
		} else if curBit.AddIfPossible(v) {
			continue
		} else {
			curBit = c.AddBit(v)
			continue
		}
	}

	return nil
}

func (c *compressor) toString() string {
	res := make([]string, 0)
	for i, _ := range c.compressed {
		res = append(res, strconv.Itoa(c.compressed[i].start))
		res = append(res, strconv.Itoa(c.compressed[i].GetOffset()))
	}
	return strings.Join(res, " ")
}

/*
Считывание данных
*/

// readCountDataSet - считывание количество наборов данных
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

// readDataLength - считывание длины набора данных
func readDataLength(r *bufio.Reader, min int, max int) (int, error) {
	cntStr, _ := r.ReadString('\n')
	if Debug {
		println("data length", cntStr)
	}
	cntStr = strings.TrimSpace(cntStr)
	cnt, err := strconv.Atoi(cntStr)
	if err != nil || cnt < min || cnt > max {
		return 0, errors.New("incorrect data length")
	}
	return cnt, nil
}

// readLine - чтение строки с данными
func readLine(r *bufio.Reader, min int, max int) ([]int, error) {
	lBytes, _, err := r.ReadLine()
	l := string(lBytes)
	if Debug {
		println("line", l)
	}
	if err != nil {
		return nil, fmt.Errorf("incorrect input: %v", err.Error())
	}
	l = strings.TrimSpace(l)
	dataStr := strings.Split(l, " ")

	res := make([]int, len(dataStr))
	for i, _ := range dataStr {
		res[i], err = strconv.Atoi(dataStr[i])
		if err != nil {
			return nil, err
		}
		if res[i] < min || res[i] > max {
			return nil, errors.New("incorrect data")
		}
	}

	return res, nil
}
