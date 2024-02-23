package testcase

import (
	"bufio"
	"bytes"
	"io/fs"
	"log"
	"os"
	"strings"
)

type TestCase struct {
	Name   string
	expect string
	input  string

	outBuff *bytes.Buffer
	in      *bufio.Reader
	out     *bufio.Writer
}

func (t *TestCase) GetActual() string {
	if err := t.out.Flush(); err != nil {
		log.Fatal(err.Error())
	}
	return strings.TrimSpace(t.outBuff.String())
}

func (t *TestCase) GetExpect() string {
	return t.expect
}

func (t *TestCase) Input() *bufio.Reader {
	if t.in == nil {
		input := bytes.NewBuffer([]byte(t.input))
		t.in = bufio.NewReader(input)
	}
	return t.in
}

func (t *TestCase) Output() *bufio.Writer {
	if t.out == nil {
		t.outBuff = bytes.NewBuffer(make([]byte, 0))
		t.out = bufio.NewWriter(t.outBuff)
	}
	return t.out
}

func ReadTestCase(dir string) (map[string]*TestCase, error) {
	fileSystem := os.DirFS(dir)

	testCases := make(map[string]*TestCase)

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if path == "." {
			return nil
		}

		tsName := strings.TrimSuffix(path, ".a")
		var ts *TestCase
		var ok bool
		if ts, ok = testCases[tsName]; !ok {
			ts = &TestCase{Name: tsName}
			testCases[tsName] = ts
		}

		fileData, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(path, ".a") {
			ts.expect = strings.TrimSpace(string(fileData))
		} else {
			ts.input = string(fileData)
		}
		return nil
	})

	return testCases, err
}
