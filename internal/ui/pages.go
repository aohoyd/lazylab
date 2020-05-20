package ui

import (
	"github.com/rivo/tview"
)

type pages struct {
	*tview.Pages

	active  string
	mapping map[string]tview.Primitive
}

func newPages() *pages {
	return &pages{tview.NewPages(), "", make(map[string]tview.Primitive)}
}

// AddPage adds page with given name to pages. If name is empty it will generate it
// as count_of_pages
func (mp *pages) AddPage(name string, page tview.Primitive) *pages {
	if name == "" {
		return nil
	}
	mp.Pages.AddPage(name, page, true, false)
	mp.mapping[name] = page
	return mp
}

func (mp *pages) SwitchToPage(name string) tview.Primitive {
	el, ok := mp.mapping[name]
	if !ok {
		return nil
	}
	mp.active = name
	mp.Pages.SwitchToPage(name)
	return el
}

func (mp *pages) Active() tview.Primitive {
	return mp.mapping[mp.active]
}

func (mp *pages) IsPageExists(name string) bool {
	_, ok := mp.mapping[name]
	return ok
}

func (mp *pages) Clear() {
	mp.active = ""
	for k := range mp.mapping {
		delete(mp.mapping, k)
		mp.Pages.RemovePage(k)
	}
}
