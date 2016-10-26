package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// var currentId int
var todos Todos
var session mgo.DialInfo
var todosCollection *mgo.Collection

// Give us some seed data
func init() {
	// RepoCreateTodo(Todo{Name: "Write presentation"})
	// RepoCreateTodo(Todo{Name: "Host meetup"})

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	todosCollection = session.DB("test").C("todos")

	// err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	// 	&Person{"Cla", "+55 53 8402 8510"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result := Person{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Phone:", result.Phone)
}

func FindAll() Todos {
	result := Todos{}

	err := todosCollection.Find(bson.M{}).All(&result)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("result:", result)

	return result
}

func RepoFindTodo(id bson.ObjectId) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

func RepoCreateTodo(t Todo) Todo {
	// currentId += 1
	// t.Id = currentId
	// todos = append(todos, t)
	t.Id = bson.NewObjectId()
	err := todosCollection.Insert(t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func RepoDestroyTodo(id bson.ObjectId) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
