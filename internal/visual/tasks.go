package visual

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/b0ralgin/todoist.ci/internal/tasks"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

type TasksListWidget struct {
	*gocui.View
	cli     *tasks.Client
	tasks   []tasks.Task
	size    int
	offset  int
	total   int
	current int
}

const (
	TasksViewName      = "tasks"
	ChangeTaskViewName = "change_task"
)

func NewTasksListWidget(cli *tasks.Client) *TasksListWidget {
	return &TasksListWidget{cli: cli}
}

func (tw *TasksListWidget) Layout(g *gocui.Gui) error {
	maxX, MaxY := g.Size()
	v, err := g.SetView(TasksViewName, 1, 0, maxX-1, MaxY-4)
	if err != nil {
		if errors.Is(err, gocui.ErrUnknownView) {
			tw.View = v
			tw.Autoscroll = false
			tw.SelBgColor = gocui.ColorBlack
			tw.SelFgColor = gocui.ColorWhite | gocui.AttrBold
			tw.Highlight = true
			if _, err := g.SetCurrentView(TasksViewName); err != nil && !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
		} else {
			return err
		}
	}
	tw.displayList()
	return nil
}

func (tw *TasksListWidget) Sync() error {
	if err := tw.cli.Sync(context.Background()); err != nil {
		return err
	}
	tasks, err := tw.cli.GetTasks()
	if err != nil {
		return err
	}
	tw.tasks = tasks
	sort.Slice(tw.tasks, func(i, j int) bool {
		return tw.tasks[i].DueTo.Before(tw.tasks[j].DueTo)
	})
	tw.total = len(tasks)
	return nil
}

func (tw *TasksListWidget) displayList() {
	tw.Clear()
	vx, vy := tw.Size()
	num := vx / 3 // widht divided by amoun of columns
	w := tabwriter.NewWriter(tw.View, num, 2, 3, ' ', 0)
	for i := tw.offset; i < tw.offset+vy; i++ {
		if i >= tw.total {
			break
		}
		t := tw.tasks[i]
		// Pri | project | content |  Date
		if _, err := color.New(colorPriority(t.Priority)).Fprintf(w, "%s\t%s\n",
			t.Text,
			t.DueTo.String(),
		); err != nil {
			panic(err)
		}
	}
	w.Flush()
	//tw.Title = fmt.Sprintf("%d %d", tw.offset, vy)
	tw.SetCursor(0, tw.current)
	tw.size = vy
}

func (tw *TasksListWidget) ScrollDown(g *gocui.Gui, v *gocui.View) error {
	if tw.current+1 > tw.size {
		if tw.offset+tw.size >= tw.total {
			return nil
		}
		tw.offset++
		return nil
	}
	if tw.current >= len(tw.tasks)-1 {
		return nil
	}
	tw.current += 1
	tw.displayList()
	return nil
}

func (tw *TasksListWidget) ScrollUp(g *gocui.Gui, v *gocui.View) error {
	if tw.current-1 < 0 {
		if tw.offset == 0 {
			return nil
		}
		tw.offset--
		return nil
	}
	tw.current -= 1
	tw.displayList()
	return nil
}

func (tw *TasksListWidget) EditTask(g *gocui.Gui, v *gocui.View) error {
	t := tw.tasks[tw.offset+tw.current]
	ts := NewTaskWidget(t, t.ID)
	return ts.AddTask(g, tw.cli, tw.Sync)
}

func (tw *TasksListWidget) NewTask(g *gocui.Gui, v *gocui.View) error {
	ts := NewTaskWidget(tasks.Task{}, -1)
	return ts.AddTask(g, tw.cli, tw.Sync)
}

func (tw *TasksListWidget) CompleteTask(g *gocui.Gui, v *gocui.View) error {
	t := tw.tasks[tw.offset+tw.current]
	v.Title = strconv.Itoa(t.ID)
	if err := tw.cli.CompleteTask(t.ID); err != nil {
		return err
	}
	if err := tw.Sync(); err != nil {
		return err
	}
	tw.displayList()
	return nil
}

func (tw *TasksListWidget) SetTaskDone(g *gocui.Gui, v *gocui.View) error {
	id := tw.tasks[tw.current+tw.offset].ID
	if err := tw.cli.CompleteTask(id); err != nil {
		return NewError("failed to save task", g)
	}
	if err := tw.Sync(); err != nil {
		return NewError("failed to sync", g)
	}
	return nil
}

func (tw *TasksListWidget) ChangeProject(g *gocui.Gui, v *gocui.View) error {
	view := g.CurrentView()
	if view.Name() != ProjectWidgetName {
		return errors.New("wrong view")
	}
	_, curY := view.Cursor()
	line, err := view.Line(curY)
	if err != nil {
		return err
	}
	current := tw.current + tw.offset
	tw.tasks[current].Project = line
	id := tw.tasks[current].ID
	if err := tw.cli.ChangeProject(id, line); err != nil {
		return NewError("failed to change project", g)
	}
	g.DeleteKeybindings(ProjectWidgetName)
	if err := g.SetKeybinding(ProjectWidgetName, gocui.KeyEnter, gocui.ModNone, tw.ChangeProject); err != nil {
		return err
	}
	if err := g.DeleteView(ProjectWidgetName); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(TasksViewName); err != nil {
		return err
	}
	return nil
}

func (tw *TasksListWidget) DeleteTask(g *gocui.Gui, v *gocui.View) error {
	t := tw.tasks[tw.current+tw.offset]
	if err := tw.cli.DeleteTask(t.ID); err != nil {
		return NewError("failed to delete task", g)
	}
	if err := tw.Sync(); err != nil {
		return NewError("failed to sync", g)
	}
	return nil
}

func mapPriority(p uint) string {
	switch p {
	case 1:
		return "!"
	case 2:
		return "!!"
	case 3:
		return "!!!"
	case 4:
		return "!!!!"
	default:
		return "0"
	}
}

func colorPriority(c uint) color.Attribute {
	switch c {
	case 4:
		return color.FgRed
	case 3:
		return color.FgYellow
	case 2:
		return color.FgBlue
	default:
		return color.Attribute(0)
	}
}
