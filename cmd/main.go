package main

import (
	"errors"
	"log"

	"github.com/b0ralgin/todoist.ci/internal/tasks"
	"github.com/b0ralgin/todoist.ci/internal/visual"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	g.InputEsc = true
	helpWidget := visual.NewHelpWidget()
	todoistCli := tasks.NewClient("c29bf706b08674e03e7de2488873f4d385e3df49")
	taskListWidget := visual.NewTasksListWidget(todoistCli)
	if err := taskListWidget.Sync(); err != nil {
		log.Panicln(err)
	}
	g.SetManager(helpWidget, taskListWidget)
	projectWidget := visual.NewProjectWidget(todoistCli.GetProjects())
	if _, err := g.SetCurrentView(visual.TasksViewName); !errors.Is(err, gocui.ErrUnknownView) {
		log.Panicln(err)
	}

	/*if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, switchView); err != nil {
		log.Panicln(err)
	}*/

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyArrowDown, gocui.ModNone, taskListWidget.ScrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyArrowUp, gocui.ModNone, taskListWidget.ScrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyCtrlE, gocui.ModNone, taskListWidget.EditTask); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyCtrlN, gocui.ModNone, taskListWidget.NewTask); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyCtrlR, gocui.ModNone, taskListWidget.CompleteTask); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyCtrlP, gocui.ModNone, projectWidget.ShowProject); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.ProjectWidgetName, gocui.KeyEnter, gocui.ModNone, taskListWidget.ChangeProject); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyCtrlD, gocui.ModNone, taskListWidget.CompleteTask); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(visual.TasksViewName, gocui.KeyDelete, gocui.ModNone, taskListWidget.DeleteTask); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
