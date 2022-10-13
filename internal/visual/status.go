package visual

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

// Look
// Project ------ donw task - sync
// Work ------ 3/5 - 35m ago

type HelpWidget struct {
	status string
}

func NewHelpWidget() *HelpWidget {
	return &HelpWidget{}
}

func (s *HelpWidget) Layout(g *gocui.Gui) error {
	maxX, MaxY := g.Size()
	v, err := g.SetView("help", 1, MaxY-3, maxX-1, MaxY-1)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	v.Clear()
	v.Frame = true
	v.Wrap = true
	currentView := g.CurrentView()
	if currentView == nil {
		return nil
	}
	switch currentView.Name() {
	case ChangeTaskViewName:
		fmt.Fprint(v, "Enter: Save, Esc: Cancel")
	case TasksViewName:
		fmt.Fprint(v, "UpArrow, DownArrow: up/down, C-E: Edit task,C-P: Change priority, C-C: Quit")
	default:
	}
	return nil
}

/*func (s *HelpWidget) Content(width int) string {
	taskNumber := strconv.Itoa(s.finishedTasks)
	lastSync := s.lastSync.String()
	template := []string{}
	template = append(template, taskNumber)
	sp1 := (width/2 - len(taskNumber) - len(s.project))
	spaces := makeSpaces(sp1, " ")
	template = append(template, spaces...)
	template = append(template, s.project)
	sp2 := (width/2 - len(s.project) - len(lastSync))
	template = append(template, makeSpaces(sp2, " ")...)
	template = append(template, lastSync)
	str := strings.Join(template, "")
	return str
}*/

func makeSpaces(n int, del string) []string {
	if n < 1 {
		return []string{}
	}
	res := make([]string, n)
	for i := range res {
		res[i] = del
	}
	return res
}
