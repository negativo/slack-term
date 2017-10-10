package context

import (
	"github.com/erroneousboat/termui"
	termbox "github.com/nsf/termbox-go"

	"slack-term/config"
	"slack-term/service"
	"slack-term/views"
)

const (
	CommandMode = "command"
	InsertMode  = "insert"
	SearchMode  = "search"
)

type AppContext struct {
	EventQueue chan termbox.Event
	Service    *service.SlackService
	Body       *termui.Grid
	View       *views.View
	Config     *config.Config
	Mode       string
}

// CreateAppContext creates an application context which can be passed
// and referenced througout the application
func CreateAppContext(flgConfig string) (*AppContext, error) {
	// Load config
	config, err := config.NewConfig(flgConfig)
	if err != nil {
		return nil, err
	}

	// Create Service
	svc := service.NewSlackService(config.SlackToken)

	// Create ChatView
	view := views.CreateChatView(svc)

	return &AppContext{
		EventQueue: make(chan termbox.Event, 20),
		Service:    svc,
		View:       view,
		Config:     config,
		Mode:       CommandMode,
	}, nil
}
