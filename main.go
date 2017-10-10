package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/erroneousboat/termui"
	termbox "github.com/nsf/termbox-go"

	"slack-term/context"
	"slack-term/handlers"
)

const (
	VERSION = "v0.2.3"
	USAGE   = `NAME:
    slack-term - slack client for your terminal

USAGE:
    slack-term -config [path-to-config]

VERSION:
    %s

GLOBAL OPTIONS:
   --help, -h
`
)

var (
	flgConfig string
	flgUsage  bool
)

func init() {
	// Get home dir for config file default
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Parse flags
	flag.StringVar(
		&flgConfig,
		"config",
		path.Join(usr.HomeDir, "slack-term.json"),
		"location of config file",
	)

	flag.Usage = func() {
		fmt.Printf(USAGE, VERSION)
	}

	flag.Parse()
}

func main() {
	// Start terminal user interface
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	// Create custom event stream for termui because
	// termui's one has data race conditions with its
	// event handling. We're circumventing it here until
	// it has been fixed.
	customEvtStream := &termui.EvtStream{
		Handlers: make(map[string]func(termui.Event)),
	}
	termui.DefaultEvtStream = customEvtStream

	// Create context
	ctx, err := context.CreateAppContext(flgConfig)
	if err != nil {
		termbox.Close()
		log.Println(err)
		os.Exit(0)
	}

	// Setup body
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(ctx.Config.SidebarWidth, 0, ctx.View.Channels),
			termui.NewCol(ctx.Config.MainWidth, 0, ctx.View.Chat),
		),
		termui.NewRow(
			termui.NewCol(ctx.Config.SidebarWidth, 0, ctx.View.Mode),
			termui.NewCol(ctx.Config.MainWidth, 0, ctx.View.Input),
		),
	)
	termui.Body.Align()
	termui.Render(termui.Body)

	// Set body in context
	ctx.Body = termui.Body

	// Register handlers
	handlers.RegisterEventHandlers(ctx)

	termui.Loop()
}
