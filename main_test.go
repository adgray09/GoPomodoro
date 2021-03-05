package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChooseAlarm(t *testing.T) {
	tests := []struct {
		alarm  string
		err    error
		result string
	}{
		{alarm: "alarm1", err: nil, result: "./assets/alarms/Alarm.mp3"},
		{alarm: "alarm2", err: nil, result: "./assets/alarms/Alarm2.mp3"},
	}
	for _, test := range tests {
		result := chooseAlarm(test.alarm)
		assert.Equal(t, test.result, result)
	}

}
