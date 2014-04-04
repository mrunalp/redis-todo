package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/hoisie/redis"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"sync"
)

type Task struct {
	Name string `form:"name"`
}

var taskList []Task

var client redis.Client

var mutex *sync.Mutex

func GetTasks() []Task {
	tList := make([]Task, 20)
	mutex.Lock()
	dbvals, _ := client.Lrange("testlist", 0, -1)
	mutex.Unlock()
	for _, v := range dbvals {
		tList = append(tList, Task{string(v)})
	}
	return tList
}

func AddTask(t Task) {
	mutex.Lock()
	fmt.Printf("Adding task %v\n", t)
	client.Lpush("testlist", []byte(t.Name))
	mutex.Unlock()
}

func main() {
	m := martini.Classic()

	taskList = make([]Task, 20)

	mutex = &sync.Mutex{}

	m.Use(render.Renderer(render.Options{Directory: "/root/gosrc/src/github.com/mrunalp/redis-todo/templates"}))

	m.Get("/", func() string {
		return "Hello Martini!"
	})

	m.Get("/tasks", func(r render.Render) {
		r.HTML(200, "list", GetTasks())
	})

	m.Post("/tasks", binding.Form(Task{}), func(task Task, r render.Render) {
		AddTask(task)
		r.HTML(200, "list", GetTasks())
	})

	m.Run()
}
