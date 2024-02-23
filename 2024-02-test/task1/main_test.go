package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"ozon-contest/pkg/testcase"
	"testing"
)

func TestStartTask(t *testing.T) {
	Debug = false

	testCases, err := testcase.ReadTestCase("./test")
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			err = startTask(tc.Input(), tc.Output())
			if err != nil {
				t.Error(err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tc.GetExpect(), tc.GetActual())
		})
	}
}
