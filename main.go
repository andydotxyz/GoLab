package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//go:embed calendar.svg
var calendarIcon []byte

func main() {
	a := app.New()
	loadTheme(a)
	w := a.NewWindow("GoLab")

	g := newGUI()
	w.SetContent(g.makeUI())

	g.setupActions()
	g.loadData()
	w.ShowAndRun()
}

// here you can add some button / callbacks code using widget IDs
func (g *gui) setupActions() {
	g.tabs.SetTabLocation(container.TabLocationBottom)
	g.tabs.Items[1].Icon = fyne.NewStaticResource("calendar.svg", calendarIcon)

	g.brandintro.ParseMarkdown(`# GOLAB

## The International Conference on Go
## in Florence`)

	for _, s := range g.brandintro.Segments {
		if t, ok := s.(*widget.TextSegment); ok {
			t.Style.Alignment = fyne.TextAlignCenter
		}
	}
}

func (g *gui) showAbout() {
	message := `This app is created by the Fyne toolkit and the Fysion app from Fyne Labs.

It is based on the work of Paul Thiel at the GoLab hackathon dinner in 2023.

Many thanks to Paul, the GoLab team and the Fyne community for making this possible.`
	dialog.ShowInformation("About GoLab App", message, fyne.CurrentApp().Driver().AllWindows()[0])
}

func (g *gui) loadData() {
	store := fyne.CurrentApp().Preferences().StringList("favourites")
	faves := binding.BindStringList(&store)
	faves.AddListener(binding.NewDataListener(func() {
		fyne.CurrentApp().Preferences().SetStringList("favourites", store)
		g.plan.Refresh()
	}))

	goLabData, err := getData("https://golab.io/schedule")
	if err != nil {
		panic(err)
	}

	g.schedule.Items = nil
	for _, d := range goLabData {
		// TODO to list
		list := container.NewVBox()
		for _, r := range d.Schedule {
			id := r.Id
			check := widget.NewCheck(recordSummary(r), func(on bool) {
				if on {
					faves.Append(id)
				} else {
					faves.Remove(id)
				}
			})
			check.Checked = contains(store, r.Id)

			list.Add(check)
		}

		g.schedule.Append(widget.NewAccordionItem(d.Title, container.NewScroll(list)))
	}

	g.plan.Length = func() int {
		return len(store)
	}
	g.plan.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		id := store[i]
		for _, d := range goLabData {
			for _, r := range d.Schedule {
				if r.Id == id {
					o.(*widget.Label).SetText(recordSummary(r))
					return
				}
			}
		}
	}
}

func contains(list []string, item string) bool {
	for _, str := range list {
		if item == str {
			return true
		}
	}

	return false
}

func recordSummary(r Record) string {
	return fmt.Sprintf("%s-%s %s", timeFormat(r.Time), timeFormat(r.Time.Add(time.Minute*time.Duration(r.DurationInMinutes))), r.Title)
}

func timeFormat(t time.Time) string {
	return t.Format("15:04")
}
