package views

import (
	"github.com/erroneousboat/termui"

	"slack-term/components"
	"slack-term/service"
)

type View struct {
	Input    *components.Input
	Chat     *components.Chat
	Channels *components.Channels
	Mode     *components.Mode
}

func CreateChatView(svc *service.SlackService) *View {
	input := components.CreateInput()

	channels := components.CreateChannels(svc, input.Par.Height)

	chat := components.CreateChat(
		svc,
		input.Par.Height,
		svc.SlackChannels[channels.SelectedChannel],
		svc.Channels[channels.SelectedChannel],
	)

	mode := components.CreateMode()

	view := &View{
		Input:    input,
		Channels: channels,
		Chat:     chat,
		Mode:     mode,
	}

	return view
}

func (v *View) Refresh() {
	termui.Render(
		v.Input,
		v.Channels,
		v.Mode,
	)
}
