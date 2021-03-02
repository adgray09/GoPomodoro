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

	workDuration := flag.Duration("work", 25*time.Minute, "Work duration - default: 25 minutes")
	breakDuration := flag.Duration("break", 5*time.Minute, "Break duration - default: 5 minutes")
	flag.Parse()

	workTimer(*workDuration)
	breakTimer(*breakDuration)

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

func createTimer(timerDuration time.Duration, action func()) *time.Timer {
	timer := time.NewTimer(timerDuration)

	go func() {
		<-timer.C
		action()
	}()

	return timer
}

func workTimer(timerDuration time.Duration) {
	timer := createTimer(timerDuration, func() {
	})

	defer func() {
		select {
		case <-timer.C:
		default:
			alarm()
			fmt.Println("Take a break! Start your break timer.")
			// newDuration := 5 * time.Second
			// breakTimer(newDuration)
		}
	}()
	<-timer.C

}

func breakTimer(timerDuration time.Duration) {
	timer := createTimer(timerDuration, func() {
	})

	defer func() {
		timer.Stop()
		select {
		case <-timer.C:
		default:
			alarm()
			fmt.Println("Break over! start work timer!")
		}
	}()
	<-timer.C
}

// Future implementation ?
// func testingTicker() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	ticker2 := time.NewTicker(30 * time.Second)
// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case t := <-ticker.C:
// 				fmt.Println("Tick at", t)
// 				alarm()
// 			case t2 := <-ticker2.C:
// 				fmt.Println("Ticker 2", t2)
// 			}
// 		}
// 	}()

// 	time.Sleep(1 * time.Hour)
// 	ticker.Stop()
// 	done <- true
// 	fmt.Println("Ticker stopped")
// }

func showTimeLeft(timerDuration time.Duration) {
	// not currently implemented
	now := time.Now()
	end := now.Add(timerDuration)

	for range time.Tick(1 * time.Second) {
		newNow := time.Now()

		dif := end.Sub(newNow)

		total := int(dif.Seconds())

		minutes := int(total/60) % 60
		seconds := int(total % 60)

		fmt.Println(minutes, seconds)
	}
}
