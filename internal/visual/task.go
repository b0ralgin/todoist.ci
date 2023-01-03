package visual

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/b0ralgin/todoist.ci/internal/tasks"
	"github.com/jroimartin/gocui"
)

type TaskWidget struct {
	selectedField string
	idx           string
	priority      uint
	date          civil.Date
	text          string
}

func NewTaskWidget(task tasks.Task, idx string) *TaskWidget {
	if idx == "" {
		return &TaskWidget{priority: 4, date: civil.DateOf(time.Now()), text: ""}
	}
	return &TaskWidget{idx: idx, priority: task.Priority, text: task.Text, date: task.DueTo}
}

func (ts *TaskWidget) AddTask(g *gocui.Gui, todoCli *tasks.Client, syncFn func() error) error {
	maxX, maxY := g.Size()
	t, err := g.SetView("add_task", maxX/4, maxY/2-2, 3*maxX/4, maxY/2)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	t.Editable = true
	t.Title = "name"
	t.Highlight = true
	fmt.Fprint(t, ts.text)
	g.Cursor = true
	p, err := g.SetView("priority", maxX/4, maxY/2+1, 2*maxX/4, maxY/2+3)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	p.Title = "priority"
	fmt.Fprint(p, mapPriority(ts.priority))
	d, err := g.SetView("date", 2*maxX/4, maxY/2+1, 3*maxX/4, maxY/2+3)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	d.Title = "dute to"
	fmt.Fprint(d, ts.date.String())
	if _, err = g.SetCurrentView("add_task"); err != nil {
		return err
	}
	ts.selectedField = "add_task"
	if err := ts.keybindings(g); err != nil {
		return err
	}
	applyFn := func(g *gocui.Gui, v *gocui.View) error {
		defer func() {
			ts.Discard(g, v)
		}()
		if err := todoCli.EditTask(tasks.Task{
			ID:       ts.idx,
			Priority: ts.priority,
			Text:     t.Buffer(),
			DueTo:    ts.date,
		}); err != nil {
			//TODO: show error on screen
			return NewError("failed to save task", g)
		}
		if err := syncFn(); err != nil {
			return NewError("failed to sync", g)
		}
		return nil
	}

	if err := g.SetKeybinding("priority", gocui.KeyCtrlS, gocui.ModNone, applyFn); err != nil {
		return err
	}
	if err := g.SetKeybinding("date", gocui.KeyCtrlS, gocui.ModNone, applyFn); err != nil {
		return err
	}
	if err := g.SetKeybinding("add_task", gocui.KeyCtrlS, gocui.ModNone, applyFn); err != nil {
		return err
	}
	return nil
}

func (ts *TaskWidget) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ts.ChangeField); err != nil {
		return err
	}
	if err := g.SetKeybinding("priority", gocui.KeyArrowUp, gocui.ModNone, ts.IncPriority); err != nil {
		return err
	}
	if err := g.SetKeybinding("priority", gocui.KeyArrowDown, gocui.ModNone, ts.DecPriority); err != nil {
		return err
	}
	if err := g.SetKeybinding("date", gocui.KeyArrowDown, gocui.ModNone, ts.DueTo(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("date", gocui.KeyArrowUp, gocui.ModNone, ts.DueTo(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("date", gocui.KeyCtrlW, gocui.ModNone, ts.DueTo(7)); err != nil {
		return err
	}
	if err := g.SetKeybinding("date", gocui.KeyCtrlUnderscore, gocui.ModNone, ts.RemoveDueTo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, ts.Discard); err != nil {
		return err
	}

	return nil
}

func (ts *TaskWidget) DueTo(days int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		t := ts.date.AddDays(days)
		if ct := civil.DateOf(time.Now()); t.Before(ct) {
			ts.date = ct
			return nil
		}
		ts.date = ts.date.AddDays(days)
		v.Clear()
		fmt.Fprint(v, ts.date.String())
		return nil
	}
}

func (ts *TaskWidget) RemoveDueTo(g *gocui.Gui, v *gocui.View) error {
	ts.date = civil.Date{}
	v.Clear()
	fmt.Fprint(v, ts.date.String())
	return nil
}

func (ts *TaskWidget) ChangeField(g *gocui.Gui, v *gocui.View) error {
	switch ts.selectedField {
	case "add_task":
		ts.selectedField = "priority"
	case "priority":
		ts.selectedField = "date"
	case "date":
		ts.selectedField = "add_task"
	}
	if _, err := g.SetCurrentView(ts.selectedField); err != nil {
		return err
	}
	return nil
}

func (ts *TaskWidget) IncPriority(g *gocui.Gui, v *gocui.View) error {
	if ts.priority == 4 {
		return nil
	}
	ts.priority++
	v.Clear()
	fmt.Fprint(v, mapPriority(ts.priority))
	return nil
}

func (ts *TaskWidget) DecPriority(g *gocui.Gui, v *gocui.View) error {
	if ts.priority == 0 {
		return nil
	}
	ts.priority--
	v.Clear()
	fmt.Fprint(v, mapPriority(ts.priority))
	return nil
}

func (ts *TaskWidget) Discard(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("add_task")
	g.DeleteKeybindings("add_task")
	g.DeleteView("priority")
	g.DeleteKeybindings("priority")
	g.DeleteView("date")
	g.DeleteKeybindings("date")
	g.SetCurrentView(TasksViewName)
	return nil

}
