package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

type Question struct {
	Text  string `form:"text"`
	Score int    `form:"score"`
}

// DB Returns a martini.Handler
func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("advent"))
		defer s.Close()
		c.Next()
	}
}

// GetAll returns all Questions in the database
func GetAll(db *mgo.Database) []Question {
	var questions []Question
	db.C("questions").Find(nil).All(&questions)
	return questions
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB())

	m.Get("/questions", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "questions", GetAll(db))
	})

	m.Post("/questions", binding.Form(Question{}), func(question Question, r render.Render, db *mgo.Database) {
		db.C("questions").Insert(question)
		r.HTML(200, "questions", GetAll(db))
	})

	m.Run()
}
