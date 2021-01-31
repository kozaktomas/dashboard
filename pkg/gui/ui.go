package gui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kozaktomas/dashboard/pkg/gitlab"
	"github.com/kozaktomas/dashboard/pkg/utils"
	"log"
	"time"
)

const helpTimeout = 2 * time.Second
const helpDefault = " 💡 [o / Enter] - open in browser   [Tab / 1 / 2] - switch tab]   [b] - copy branch name"

type tabItem struct {
	text string
	url  string
	copy string
}

type tab struct {
	title  string
	number int
	label  string
	items  []tabItem
	window *widgets.List
}

type gui struct {
	width  int
	height int
	tabs   []*tab
	flash  string
	gitlab *gitlab.Service
}

func New(g *gitlab.Service) *gui {
	return &gui{
		width:  0,
		height: 0,
		flash:  "",
		tabs:   []*tab{},
		gitlab: g,
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

	mrs, err := g.gitlab.GetMyMergeRequests()
	if err != nil {
		g.flash = err.Error()
	}
	var myMrItems []tabItem
	for _, mr := range mrs {
		myMrItems = append(myMrItems, tabItem{
			text: formatGitlabMergeRequestTitle(mr),
			url:  mr.Url,
			copy: mr.BranchName,
		})
	}
	mrTab := createTab("My Merge Requests", 0, myMrItems)

	mrs, err = g.gitlab.GetReviews()
	if err != nil {
		g.flash = err.Error()
	}
	var reviews []tabItem
	for _, mr := range mrs {
		reviews = append(reviews, tabItem{
			text: formatGitlabMergeRequestTitle(mr),
			url:  mr.Url,
			copy: mr.BranchName,
		})
	}
	crTab := createTab("Code reviews", 1, reviews)

	g.tabs = append(g.tabs, mrTab)
	g.tabs = append(g.tabs, crTab)

	g.render()

	g.tabs[0].window.TitleStyle.Bg = ui.ColorYellow // first tab is active
	ui.Render(g.tabs[0].window)                     // rerender first tab

	g.controlLoop()
}

func createTab(title string, number int, items []tabItem) *tab {
	return &tab{
		title:  title,
		number: number,
		label:  fmt.Sprintf("%d", number+1),
		items:  items,
	}
}

func (g *gui) controlLoop() {
	uiEvents := ui.PollEvents()
	active := 0

	for {
		e := <-uiEvents

		if e.ID == "<Up>" || e.ID == "k" {
			g.tabs[active].window.ScrollUp()
		}

		if e.ID == "<Down>" || e.ID == "j" {
			g.tabs[active].window.ScrollDown()
		}

		for _, tab := range g.tabs {
			if e.ID == tab.label {
				g.tabs[active].window.TitleStyle.Bg = ui.ColorBlack
				ui.Render(g.tabs[active].window) // render old one
				active = tab.number
				g.tabs[active].window.TitleStyle.Bg = ui.ColorYellow
				break
			}
		}

		if e.ID == "<Tab>" {
			g.tabs[active].window.TitleStyle.Bg = ui.ColorBlack
			ui.Render(g.tabs[active].window) // render old one

			c := len(g.tabs)
			s := g.tabs[active].number
			if s+1 < c {
				active++
			} else {
				active = 0
			}

			g.tabs[active].window.TitleStyle.Bg = ui.ColorYellow
		}

		if e.ID == "q" || e.ID == "<C-c>" {
			return
		}

		if e.ID == "<Enter>" || e.ID == "o" || e.ID == "O" {
			selected := g.tabs[active].window.SelectedRow
			url := g.tabs[active].items[selected].url
			utils.OpenBrowser(url)
		}

		if e.ID == "b" || e.ID == "B" {
			selected := g.tabs[active].window.SelectedRow
			branchName := g.tabs[active].items[selected].copy
			utils.CopyToClipboard("git checkout " + branchName)
			g.RenderHelp("📋 Copied! [ " + branchName + " ]")
		}

		if g.tabs[active] != nil {
			ui.Render(g.tabs[active].window)
		}
	}
}

func (g *gui) render() {
	delta := 0
	if g.flash != "" {
		delta = g.renderFlash(delta)
	}
	g.renderTabs(delta)
	g.RenderHelp("")
}

func (g *gui) renderFlash(delta int) int {
	p := widgets.NewParagraph()
	p.Text = "  " + g.flash
	p.TextStyle.Fg = ui.ColorRed
	p.SetRect(0, 0, g.width, 3+delta)
	ui.Render(p)

	return 3 + delta
}

func (g *gui) renderTabs(delta int) {
	c := len(g.tabs)
	for i, tab := range g.tabs {
		l := widgets.NewList()
		l.Title = fmt.Sprintf(" %s [%d] ", tab.title, tab.number+1)
		var rows []string
		for _, item := range tab.items {
			rows = append(rows, item.text)
		}
		l.Rows = rows
		l.TextStyle = ui.NewStyle(ui.ColorGreen)
		l.WrapText = false
		l.SetRect(i*g.width/c, 0+delta, (i+1)*g.width/c-1, g.height-3)
		tab.window = l
		ui.Render(l)
	}
}

func (g *gui) RenderHelp(info string) {
	p := widgets.NewParagraph()
	if info == "" {
		p.Text = helpDefault
	} else {
		p.Text = info
		ui.Render(p)
		go func() {
			time.Sleep(helpTimeout)
			p.Text = helpDefault
			ui.Render(p)
		}()
	}
	p.SetRect(0, g.height-3, g.width-1, g.height)
	ui.Render(p)
}
