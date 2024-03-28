package main

import (
    "fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)


var (
    logger zerolog.Logger
)

func ConsoleLogConfig(w *zerolog.ConsoleWriter) {
    w.Out = os.Stdout
    w.NoColor = false
    //w.TimeFormat = time.RFC3339
    w.TimeFormat = time.DateTime + ".000-07"
    w.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("[%-5s]", i)) }
    w.FormatTimestamp = func(i interface{}) string {return time.Now().Format(w.TimeFormat)}

}

func init() {
    out := zerolog.NewConsoleWriter(ConsoleLogConfig)
    logger = zerolog.New(out)
}
