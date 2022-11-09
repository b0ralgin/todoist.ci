package visual

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

func NewError(msg string, g *gocui.Gui) error {
	sizeX, sizeY := g.Size()
	v, err := g.SetView("error", sizeX/2-len(msg)/2, sizeY/2-1, sizeX/2+len(msg)/2, sizeY/2+1)
	if err != nil {
		return err
	}
	fmt.Fprintln(v, msg)
	time.Sleep(1 * time.Second)
	if err := g.DeleteView("error"); err != nil {
		return err
	}
	return nil
}
