/*

Copyright (C) 2017-2019  Ettore Di Giacinto <mudler@gentoo.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package service

import (
	"sort"

	"github.com/go-test/deep"

	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	citasks "github.com/MottainaiCI/mottainai-server/pkg/tasks"
	schema "github.com/MottainaiCI/mottainai-server/routes/schema"
	v1 "github.com/MottainaiCI/mottainai-server/routes/schema/v1"

	"github.com/mudler/anagent"
)

type Bridge interface {
	Listen(e, fn interface{})
	Publish(e, thing interface{})
	Run()
}

type ClientService struct {
	*anagent.Anagent
	Client   client.HttpClient
	Tracked  TrackHash
	PollTime int64
}

type TrackHash map[string]TaskMap
type TaskUpdate struct {
	ID   string
	Task TaskMap
	Diff TaskMap
}
type TaskMap map[string]interface{}

func NewBridge(c client.HttpClient) Bridge {
	return &ClientService{Client: c, Tracked: make(TrackHash), Anagent: anagent.New(), PollTime: int64(2)}
}

func (c *ClientService) Listen(e, fn interface{}) {
	c.Anagent.Emitter().On(e, fn)
}

func (c *ClientService) Publish(e, thing interface{}) {
	c.Anagent.Emitter().Emit(e, thing)
}

func (c *ClientService) Run() {
	c.TrackTasks()

	c.Map(c)
	c.TimerSeconds(c.PollTime, true, func(d *ClientService) {
		new, updated, removed := d.TrackTasks()
		for _, v := range new {
			d.Publish("task.created", v)
		}
		for id, v := range updated {
			d.Publish("task.update", &TaskUpdate{ID: id, Task: c.Tracked[id], Diff: v})
		}
		for _, v := range removed {
			d.Publish("task.removed", v)
		}
	})

	c.Start()
}

func (d *ClientService) TrackTasks() (TrackHash, map[string]TaskMap, map[string]TaskMap) {
	tasks := d.TaskList()
	new := make(TrackHash)
	updated := make(map[string]TaskMap)
	removed := make(map[string]TaskMap)

	seen := make(map[string]bool)

	for _, v := range tasks {
		seen[v.ID] = true
		current_entry := TaskMap(v.ToMap())

		if _, ok := d.Tracked[v.ID]; !ok {
			d.Tracked[v.ID] = current_entry
			d.Tracked[v.ID]["ID"] = v.ID
			new[v.ID] = current_entry
			new[v.ID]["ID"] = v.ID
		}
		current_entry["ID"] = v.ID
		old_entry := d.Tracked[v.ID]

		for field, _ := range old_entry {
			if diff := deep.Equal(current_entry[field], old_entry[field]); diff != nil {
				if len(updated[v.ID]) == 0 {
					updated[v.ID] = make(TaskMap)
				}
				updated[v.ID][field] = old_entry[field]
			}
		}

		d.Tracked[v.ID] = current_entry
		d.Tracked[v.ID]["ID"] = v.ID
	}

	for id, m := range d.Tracked {
		if s, ok := seen[id]; !ok || !s {
			removed[id] = m
			delete(d.Tracked, id)
		}
	}
	return new, updated, removed
}

func (c *ClientService) TaskList() []citasks.Task {
	var tlist []citasks.Task
	req := schema.Request{
		Route:  v1.Schema.GetTaskRoute("show_all"),
		Target: &tlist,
	}
	if err := c.Client.Handle(req); err != nil {
		return tlist
	}

	sort.Slice(tlist[:], func(i, j int) bool {
		return tlist[i].CreatedTime > tlist[j].CreatedTime
	})
	return tlist
}
