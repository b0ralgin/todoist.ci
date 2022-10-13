package visual

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type ProjectWidget struct {
	projects map[int]string
	selected string
}

const (
	ProjectWidgetName = "projects"
)

func NewProjectWidget(projects map[int]string) *ProjectWidget {
	return &ProjectWidget{projects: projects}
}

func (pw *ProjectWidget) ShowProject(g *gocui.Gui, v *gocui.View) error {
	sizeX, sizeY := g.Size()
	view, err := g.SetView(ProjectWidgetName, sizeX/2-5, sizeY/2, sizeX/2+5, sizeY/2+len(pw.projects)+1)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	view.Highlight = true
	view.Autoscroll = false
	view.SelBgColor = gocui.ColorBlack
	view.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	projects := []string{}
	for _, p := range pw.projects {
		projects = append(projects, p)
	}
	fmt.Fprint(view, strings.Join(projects, "\n"))
	if _, err := g.SetCurrentView(ProjectWidgetName); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectWidgetName, gocui.KeyArrowDown, gocui.ModNone, pw.ScrollDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(ProjectWidgetName, gocui.KeyArrowUp, gocui.ModNone, pw.ScrollUp); err != nil {
		return err
	}
	return nil
}

func (pw *ProjectWidget) GetProject(g *gocui.Gui, v *gocui.View) error {
	_, curY := v.Cursor()
	line, err := v.Line(curY)
	if err != nil {
		return err
	}
	pw.selected = line
	g.DeleteKeybindings(ProjectWidgetName)
	if err := g.DeleteView(ProjectWidgetName); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(TasksViewName); err != nil {
		return err
	}
	return nil
}

func (pw *ProjectWidget) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func (pw *ProjectWidget) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}
