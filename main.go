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

func main() {

	duration := flag.Duration("duration", 25*time.Minute, "Pomodoro duration")
	flag.Parse()

	startTimer(*duration)
}

func alarm() {
	f, err := os.Open("./assets/alarms/Alarm.mp3")
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

func timerDuration(d time.Duration) time.Duration {
	return d * time.Minute
}

func createTimer(timerDuration time.Duration, action func()) *time.Timer {
	timer := time.NewTimer(timerDuration)

	go func() {
		<-timer.C
		action()
	}()

	return timer
}

func startTimer(timerDuration time.Duration) {
	timer := createTimer(timerDuration, func() {
		fmt.Println("Take a break!")
		alarm()
	})
	defer timer.Stop()

	killTimer := (time.Duration(timerDuration) + time.Second)
	countdownBeforeExit := time.NewTimer(time.Second * killTimer)
	<-countdownBeforeExit.C
}

func minuteToSeconds(timeDuration time.Duration) int {

	convertedTime := int(timeDuration / time.Second)

	return convertedTime

}
