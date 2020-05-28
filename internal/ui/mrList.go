package ui

import (
	"fmt"
	"github.com/aovlllo/lazylab/internal/gitlab"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

type mrPage struct {
	*tview.Flex

	description *tview.TextView
	info        *tview.TextView
}

func newMRPage(mr *gitlab.MergeRequest) *mrPage {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	description := tview.NewTextView().SetText(mr.Description)
	description.SetBorder(true).SetTitle(mr.Title)

	color := "green"
	if mr.HasConflicts {
		color = "red"
	}
	infoText := fmt.Sprintf(
		"URL: %s\nAuthor: [green]%s[white]\nReference: %s\n[%s]Has conflicts: %t[white]",
		mr.WebURL, mr.Author.Name, mr.References.Full, color, mr.HasConflicts)
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText(infoText).
		SetWrap(false)
	info.SetBorder(true)

	page.AddItem(description, 0, 1, false).
		AddItem(info, 7, 1, false)

	return &mrPage{page, description, info}
}

type mrList struct {
	*tview.Flex

	app   *Application
	list  *tview.List
	pages *pages

	initialized bool
	refresh     func()
}

func NewMRList(app *Application) *mrList {
	pages := newPages()
	pages.SetBorder(true).SetTitle("Loading...")
	list := tview.NewList().SetMainTextColor(tcell.ColorYellow).SetHighlightFullLine(true)
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		pages.SwitchToPage(strconv.Itoa(index))
	})
	main := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(list, 0, 1, true).
		AddItem(pages, 0, 2, false)

	ml := &mrList{main, app, list, pages, false, nil}

	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEsc:
			app.hideFullScreen()
			app.SetFocus(list)
			return nil
		}
		return event
	})
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Rune() == 'r':
			ml.Refresh()
		case event.Rune() == 'j':
			return tcell.NewEventKey(tcell.KeyDown, rune(tcell.KeyDown), tcell.ModNone)
		case event.Rune() == 'k':
			return tcell.NewEventKey(tcell.KeyUp, rune(tcell.KeyUp), tcell.ModNone)
		case event.Key() == tcell.KeyEnter:
			app.showFullScreen()
			active := pages.Active().(*mrPage)
			if active != nil {
				app.SetFocus(active.description)
			}
		}
		return event
	})
	list.SetBorder(true).SetTitle("MRs")
	return ml
}

func (m *mrList) AddMR(mr *gitlab.MergeRequest) {
	page := newMRPage(mr)
	name := strconv.Itoa(len(m.pages.mapping))
	m.pages.AddPage(name, page)

	var daysPhrase string
	days := int(time.Since(*mr.UpdatedAt).Hours()) / 24
	switch days {
	case 0:
		daysPhrase = "updated today"
	case 1:
		daysPhrase = fmt.Sprintf("updated %d day ago", days)
	default:
		daysPhrase = fmt.Sprintf("updated %d days ago", days)
	}

	m.list.AddItem(tview.Escape(mr.Title), daysPhrase, 0, func() { m.pages.SwitchToPage(name) })
}

func (m *mrList) Clear() {
	m.pages.Clear()
	m.list.Clear()
}

func (m *mrList) Refresh() {
	m.initialized = true
	m.Clear()
	m.pages.SetBorder(true)
	m.app.ForceDraw()
	if m.refresh != nil {
		m.refresh()
	}
	m.pages.SetBorder(false)
	m.app.ForceDraw()
}

func (m *mrList) SetRefresh(fn func()) {
	m.refresh = fn
}

func (m *mrList) toggleFocus() {
	switch {
	case m.list.HasFocus():
		m.app.SetFocus(m.pages)
	case m.pages.HasFocus():
		m.app.SetFocus(m.list)
	}
}
