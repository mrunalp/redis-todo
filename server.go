package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/hoisie/redis"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

type Task struct {
	Name string `form:"name"`
}

type AddTaskOp struct {
	task Task
	resp chan bool
}

var taskList []Task

var client redis.Client

func GetTasks() []Task {
	tList := make([]Task, 20)
	dbvals, _ := client.Lrange("testlist", 0, -1)
	for _, v := range dbvals {
		tList = append(tList, Task{string(v)})
	}
	return tList
}

func AddTaskProc(c chan AddTaskOp) {
	for {
		aop := <-c
		fmt.Println("Adding task %v", aop.task)
		taskList = append(taskList, aop.task)
		client.Lpush("testlist", []byte(aop.task.Name))
		aop.resp <- true
	}
}

func AddTask(t Task, c chan AddTaskOp) {
	done := make(chan bool)
	aop := AddTaskOp{t, done}
	c <- aop
	<-done
}

func main() {
	m := martini.Classic()

	taskList = make([]Task, 20)

	ch := make(chan AddTaskOp)

	go AddTaskProc(ch)

	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "Hello Martini!"
	})

	m.Get("/tasks", func(r render.Render) {
		r.HTML(200, "list", GetTasks())
	})

	m.Post("/tasks", binding.Form(Task{}), func(task Task, r render.Render) {
		AddTask(task, ch)
		r.HTML(200, "list", GetTasks())
	})

	m.Run()
}
