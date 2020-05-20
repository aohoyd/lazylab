package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type ApplicationPrimitive interface {
	SetFocus(p tview.Primitive) *tview.Application
	Draw() *tview.Application
	QueueUpdateDraw(func()) *tview.Application
}

// Application is our application with base menu and pages on it
type Application struct {
	*tview.Application

	layout *tview.Flex
	menu   *tview.TextView
	pages  *pages
}

func newMenu(app *Application) *tview.TextView {
	menu := tview.NewTextView().
		SetDynamicColors(false).
		SetTextColor(tcell.ColorDarkCyan).
		SetRegions(true).
		SetWrap(false)

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Rune() == 'j' || event.Rune() == 'k' || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown:
			fallthrough
		case event.Key() == tcell.KeyEnter:
			active := app.pages.Active()
			if active != nil {
				app.SetFocus(active)
			}
		case event.Key() == tcell.KeyRight || event.Rune() == 'l':
			app.highlightTransform(func(i, c int) int { return (i + 1) % c })
		case event.Key() == tcell.KeyLeft || event.Rune() == 'h':
			app.highlightTransform(func(i, c int) int { return (i - 1 + c) % c })
		}
		return event
	})

	return menu
}

func NewApplication() *Application {
	app := &Application{Application: tview.NewApplication()}
	pages := newPages()
	menu := newMenu(app)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.AddItem(pages, 0, 1, false).
		AddItem(menu, 1, 0, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			switch {
			case !menu.HasFocus() && strings.HasSuffix(pages.active, "fullscreen"):
				app.hideFullScreen()
			default:
				app.SetFocus(menu)
			}
		case tcell.KeyTab:
			active := app.pages.Active().(*mrList)
			if active != nil {
				active.toggleFocus()
			}
			return nil
		}
		return event
	})
	app.SetRoot(layout, true)

	app.layout = layout
	app.menu = menu
	app.pages = pages

	return app
}

// GetSectionsCount returns count of sections.
// Every section have 2 pages: with full mr and with only it's description.
func (app *Application) GetSectionsCount() int {
	return len(app.pages.mapping) / 2
}

func (app *Application) AddSection(title string, item *mrList) {
	id := strconv.Itoa(app.GetSectionsCount())
	app.pages.AddPage(id, item).
		AddPage(id+"fullscreen", item.pages)
	fmt.Fprintf(app.menu, `["%s"]%s[""] `, id, title)
}

func (app *Application) Run() error {
	app.switchToTab("0")
	return app.Application.Run()
}

func (app *Application) switchToTab(name string) *mrList {
	if !app.pages.IsPageExists(name) {
		return nil
	}
	el := app.pages.SwitchToPage(name).(*mrList)
	if el != nil && !el.initialized {
		el.Refresh()
	}
	app.menu.Highlight(name)
	return el
}

func (app *Application) highlightTransform(transform func(int, int) int) {
	count := app.GetSectionsCount()
	if count < 2 {
		return
	}
	if app.pages.active == "" {
		app.switchToTab("0")
		return
	}
	highlighted, _ := strconv.Atoi(app.pages.active)
	next := strconv.Itoa(transform(highlighted, count))
	app.switchToTab(next)
}

func (app *Application) showFullScreen() {
	app.pages.SwitchToPage(app.pages.active + "fullscreen")
}

func (app *Application) hideFullScreen() {
	app.pages.SwitchToPage(strings.TrimSuffix(app.pages.active, "fullscreen"))
	app.SetFocus(app.pages.Active())
}
