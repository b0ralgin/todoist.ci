package tasks

import (
	"context"
	"time"

	"cloud.google.com/go/civil"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	todoist "github.com/sachaos/todoist/lib"
)

type Client struct {
	cli *todoist.Client
}

type Task struct {
	ID       int
	Project  string
	Priority uint
	Text     string
	DueTo    civil.Date
}

func NewClient(token string) *Client {
	cli := todoist.NewClient(&todoist.Config{
		AccessToken: token,
		//DebugMode:   true,
	})
	return &Client{cli: cli}
}

func (c *Client) Sync(ctx context.Context) error {
	return c.cli.Sync(ctx)
}

func (c *Client) GetProjects() map[int]string {
	res := map[int]string{}
	for _, p := range c.cli.Store.Projects {
		res[p.ID] = p.Name
	}
	return res
}

func (c *Client) GetTasks() ([]Task, error) {
	item := c.cli.Store.RootItem
	if item == nil {
		return nil, errors.New("empty list")
	}
	return c.traversal(item), nil
}

func (c *Client) ChangeProject(id int, project string) error {
	pid := c.cli.Store.Projects.GetIDByName(project)
	if pid == 0 {
		return errors.New("project not found ")
	}
	item := todoist.Item{}
	item.ID = id
	if err := c.cli.MoveItem(context.Background(), &item, pid); err != nil {
		return err
	}
	return nil
}

func (c *Client) EditTask(task Task) error {
	t := todoist.Item{
		BaseItem: todoist.BaseItem{
			Content: task.Text,
		},
		Priority: int(task.Priority),
	}
	t.DateString = task.DueTo.String()
	if task.DueTo.IsZero() {
		t.DateString = "null"
	}
	ctx := context.Background()
	if task.ID == 0 {
		return c.cli.AddItem(ctx, t)
	}
	t.ID = task.ID
	return c.cli.UpdateItem(ctx, t)
}

func (c *Client) CompleteTask(id int) error {
	return c.cli.CloseItem(context.Background(), []int{id})
}

func (c *Client) DeleteTask(id int) error {
	return c.cli.DeleteItem(context.Background(), []int{id})
}

func (c *Client) MapColor() map[string]color.Attribute {
	res := map[string]color.Attribute{}
	for _, p := range c.cli.Store.Projects {
		res[p.Name] = color.Attribute(p.Color - 9)
	}
	return res
}

func (c *Client) traversal(item *todoist.Item) []Task {
	if item == nil {
		return []Task{}
	}
	task := []Task{}
	if item.Due != nil && item.Checked == 0 && !item.DateTime().After(time.Now()) {
		task = []Task{c.mapTask(item)}
	}
	child := c.traversal(item.ChildItem)
	neigbor := c.traversal(item.BrotherItem)
	tl := append(child, neigbor...)
	return append(task, tl...)
}

func (c *Client) mapTask(item *todoist.Item) Task {
	task := Task{
		ID:       item.ID,
		Priority: uint(item.Priority),
		Text:     item.Content,
	}
	project := c.cli.Store.FindProject(item.ProjectID)
	if project != nil {
		task.Project = project.Name
	}
	if item.Due != nil {
		task.DueTo = civil.DateOf(item.DateTime())
	}
	return task
}
