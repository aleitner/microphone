package main

import (
	"fmt"
	"log"
	"microphone"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start reading microphone input",
				Action: func(c *cli.Context) error {

					fmt.Println("Recording. Press Ctrl-C to stop.")

					err := microphone.Init()
					if err != nil {
						log.Fatal(err)
					}
					defer microphone.Terminate()

					stream, format, err := microphone.OpenDefaultStream(44100, 2)
					if err != nil {
						log.Fatal(err)
					}

					defer stream.Close()

					stream.Start()
					defer stream.Stop()

					speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

					done := make(chan bool)
					speaker.Play(beep.Seq(stream, beep.Callback(func() {
						done <- true
					})))

					<-done
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
