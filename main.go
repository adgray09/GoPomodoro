package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type flags struct {
	isUsed bool
}

func main() {
	flagHandler()
}

func flagHandler() {
	var userAlarm string

	flag.StringVar(&userAlarm, "alarm", "alarm1", "Choices: alarm1, alarm2")
	workDuration := flag.Duration("work", 25*time.Minute, "Work duration - default: 25 minutes")
	breakDuration := flag.Duration("break", 5*time.Minute, "Break duration - default: 5 minutes")
	flag.Parse()

	firstFlag := flags{isFlagPassed("work")}
	secondFlag := flags{isFlagPassed("break")}

	userAlarm = chooseAlarm(userAlarm)

	if firstFlag.isUsed {
		fmt.Println("Starting work")
		workTimer(*workDuration, userAlarm)
	}

	if secondFlag.isUsed {
		fmt.Println("Starting break")
		breakTimer(*breakDuration, userAlarm)
	}

}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func chooseAlarm(userChoice string) string {
	var theAlarm string

	alarms := map[string]string{
		"alarm1": "./assets/alarms/Alarm.mp3",
		"alarm2": "./assets/alarms/Alarm2.mp3",
	}

	if key, ok := alarms[userChoice]; ok {
		theAlarm = key
	}
	return theAlarm
}

func alarm(alarm string) {
	f, err := os.Open(alarm)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	sr := format.SampleRate * 2
	speaker.Init(sr, sr.N(time.Second/10))

	resampled := beep.Resample(4, format.SampleRate, sr, streamer)
	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func createTimer(timerDuration time.Duration, action func()) *time.Timer {
	timer := time.NewTimer(timerDuration)

	go func() {
		<-timer.C
		action()
	}()

	return timer
}

func workTimer(timerDuration time.Duration, alarmChoice string) {
	timer := createTimer(timerDuration, func() {
	})

	defer func() {
		select {
		case <-timer.C:
		default:
			alarm(alarmChoice)
			fmt.Println("Take a break! Start your break timer.")
		}
	}()
	<-timer.C

}

func breakTimer(timerDuration time.Duration, alarmChoice string) {
	timer := createTimer(timerDuration, func() {
	})

	defer func() {
		select {
		case <-timer.C:
		default:
			// showTimeLeft(timerDuration)
			alarm(alarmChoice)
			fmt.Println("Break is over! start your work timer")
		}
	}()
	<-timer.C

}

func showTimeLeft(timerDuration time.Duration) {
	// not currently implemented
	now := time.Now()
	end := now.Add(timerDuration + (2 * time.Second))
	// fmt.Println(end)

	for range time.Tick(1 * time.Second) {
		newNow := time.Now()

		dif := end.Sub(newNow)

		total := int(dif.Seconds())

		minutes := int(total/60) % 60
		seconds := int(total % 60)

		if seconds == 0 {
			break
		}

		fmt.Printf("M: %#v, S: %#v", minutes, seconds)
	}
}
