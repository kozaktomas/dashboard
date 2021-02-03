package gui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kozaktomas/dashboard/pkg/integrations"
	"github.com/kozaktomas/dashboard/pkg/utils"
	"log"
	"time"
)

const helpTimeout = 5 * time.Second
const helpDefault = "💡 [o / Enter] - open in browser   [Tab / 1 / 2] - switch tab]   [b] - copy branch name"

type gui struct {
	width  int
	height int
	apps   []integrations.Integration

	list   *widgets.List
	detail *widgets.Paragraph
	flash  *widgets.Paragraph

	activeAppIndex int
	activeListItem int
}

func New(integrations []integrations.Integration) *gui {
	return &gui{
		width:  0,
		height: 0,
		apps:   integrations,

		list:   widgets.NewList(),
		detail: widgets.NewParagraph(),
		flash:  widgets.NewParagraph(),

		activeAppIndex: 0,
		activeListItem: 0,
	}
}

func (g *gui) Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	x, y := ui.TerminalDimensions()
	g.width = x
	g.height = y

	g.detail.SetRect(g.width/2+2, 0, g.width-2, g.height-1)
	g.list.SetRect(0, 0, g.width/2, g.height-1)
	g.flash.SetRect(0, g.height-1, g.width-1, g.height)
	g.flash.Border = false

	g.renderList()

	g.controlLoop()
}

func (g *gui) renderList() {
	app := g.apps[g.activeAppIndex]
	items := app.GetItems()

	var rows []string
	for _, item := range items {
		rows = append(rows, item.Text)
	}

	g.list.Title = app.GetName()

	g.list.Rows = rows
	g.list.TextStyle = ui.NewStyle(ui.ColorGreen)
	g.list.WrapText = false

	ui.Render(g.list)
	g.renderDetail()
	g.renderFlash("")
}

func (g *gui) renderDetail() {
	item := g.apps[g.activeAppIndex].GetItems()[g.list.SelectedRow]
	itemId := item.Id

	go func() {
		time.Sleep(500 * time.Millisecond)
		if itemId != g.apps[g.activeAppIndex].GetItems()[g.list.SelectedRow].Id {
			return // selected changed
		}

		detail := g.apps[g.activeAppIndex].GetDetail(item)
		text := ""
		for _, part := range detail.Parts {
			text += " " + part.GetText()
		}

		g.detail.Text = text
		ui.Render(g.detail)
	}()
}

func (g *gui) controlLoop() {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		if e.ID == "<Up>" || e.ID == "k" {
			g.list.ScrollUp()
			g.renderDetail()
		}

		if e.ID == "<Down>" || e.ID == "j" {
			g.list.ScrollDown()
			g.renderDetail()
		}

		if e.ID == "q" || e.ID == "<C-c>" {
			return
		}

		if e.ID == "<Enter>" || e.ID == "o" || e.ID == "O" {
			selected := g.list.SelectedRow
			url := g.apps[g.activeAppIndex].GetItems()[selected].Url
			utils.OpenBrowser(url)
		}

		if e.ID == "b" || e.ID == "B" {
			selected := g.list.SelectedRow
			cp := g.apps[g.activeAppIndex].GetItems()[selected].Copy
			utils.CopyToClipboard(cp)
			g.renderFlash("📋 Copied! [ " + cp + " ]")
		}

		g.renderList()
	}
}

func (g *gui) renderFlash(msg string) {
	if msg == "" {
		g.flash.Text = helpDefault
	} else {
		g.flash.Text = msg
		ui.Render(g.flash)
		go func() {
			time.Sleep(helpTimeout)
			g.flash.Text = helpDefault
			ui.Render(g.flash)
		}()
	}

	ui.Render(g.flash)
}

func (g *gui) processError(err error) {
	if err != nil {
		g.renderFlash(err.Error())
	}
}
