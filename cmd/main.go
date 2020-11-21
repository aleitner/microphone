package main

import (
	"fmt"
	"log"
	"microphone"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/malgo"
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

					ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
						fmt.Printf("LOG <%v>\n", message)
					})
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					defer func() {
						_ = ctx.Uninit()
						ctx.Free()
					}()

					deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
					deviceConfig.Capture.Format = malgo.FormatS24
					deviceConfig.Capture.Channels = 2
					deviceConfig.SampleRate = 44100

					stream, format, err := microphone.OpenStream(ctx, deviceConfig)
					if err != nil {
						log.Fatal(err)
					}

					defer stream.Close()

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
