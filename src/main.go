package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const address = "127.0.0.1:3000"

type IndexPage struct {
	Todos        []Todo
	HasCompleted bool
	TodosLeft    int
	Page         string
}

type Todo struct {
	Id        int64
	Todo      string
	Completed bool
	Editing   bool
}

func main() {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		panic("BASE_PATH environment variable not set")
	}

	db, err := sql.Open("sqlite3", path.Join(basePath, "sqlite.db"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		todosLeft := 0
		for _, todo := range todos {
			if !todo.Completed {
				todosLeft += 1
			}
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "index.tmpl"),
			path.Join(basePath, "src", "todos.tmpl"),
			path.Join(basePath, "src", "todo.tmpl"),
			path.Join(basePath, "src", "clear-completed.tmpl"),
			path.Join(basePath, "src", "active-counter.tmpl"),
			path.Join(basePath, "src", "complete-all.tmpl"),
		}...))
		data := IndexPage{
			Todos:        todos,
			HasCompleted: len(todos) > todosLeft,
			TodosLeft:    todosLeft,
			Page:         "all",
		}
		tmpl.ExecuteTemplate(w, "index.tmpl", data)
	})

	r.Get("/active", func(w http.ResponseWriter, r *http.Request) {
		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		activeTodos := []Todo{}
		for _, todo := range todos {
			if !todo.Completed {
				activeTodos = append(activeTodos, todo)
			}
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "index.tmpl"),
			path.Join(basePath, "src", "todos.tmpl"),
			path.Join(basePath, "src", "todo.tmpl"),
			path.Join(basePath, "src", "clear-completed.tmpl"),
			path.Join(basePath, "src", "active-counter.tmpl"),
			path.Join(basePath, "src", "complete-all.tmpl"),
		}...))
		data := IndexPage{
			Todos:        activeTodos,
			HasCompleted: len(todos) > len(activeTodos),
			TodosLeft:    len(activeTodos),
			Page:         "active",
		}
		tmpl.ExecuteTemplate(w, "index.tmpl", data)
	})

	r.Get("/completed", func(w http.ResponseWriter, r *http.Request) {
		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		completedTodos := []Todo{}
		for _, todo := range todos {
			if todo.Completed {
				completedTodos = append(completedTodos, todo)
			}
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "index.tmpl"),
			path.Join(basePath, "src", "todos.tmpl"),
			path.Join(basePath, "src", "todo.tmpl"),
			path.Join(basePath, "src", "clear-completed.tmpl"),
			path.Join(basePath, "src", "active-counter.tmpl"),
			path.Join(basePath, "src", "complete-all.tmpl"),
		}...))
		data := IndexPage{
			Todos:        completedTodos,
			HasCompleted: len(completedTodos) > 0,
			TodosLeft:    len(todos) - len(completedTodos),
			Page:         "completed",
		}
		tmpl.ExecuteTemplate(w, "index.tmpl", data)
	})

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		type TodosPage struct {
			Todos []Todo
			Page  string
		}

		page := "all"
		if strings.Contains(r.Referer(), "active") {
			page = "active"
		} else if strings.Contains(r.Referer(), "completed") {
			page = "completed"
		}

		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		filteredTodos := []Todo{}
		for _, todo := range todos {
			if page == "all" {
				filteredTodos = append(filteredTodos, todo)
			} else if page == "active" && !todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			} else if page == "completed" && todo.Completed {
				filteredTodos = append(filteredTodos, todo)
			}
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "todos.tmpl"),
			path.Join(basePath, "src", "todo.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "todos.tmpl", TodosPage{Todos: filteredTodos, Page: page})
	})

	r.Get("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		todo, err := findTodo(id, db)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "todo.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "todo.tmpl", *todo)
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := r.PostForm
		todo := strings.Trim(form.Get("todo"), " ")

		rows, err := db.Query("insert into todos(todo) values(?) returning id, todo, completed;", todo)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		defer rows.Close()

		var id int64
		var completed bool
		rows.Next()
		if err := rows.Scan(&id, &todo, &completed); err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "newTodo")
		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "todo.tmpl"),
		}...))
		data := Todo{Id: id, Todo: todo, Completed: completed, Editing: false}
		tmpl.ExecuteTemplate(w, "todo.tmpl", data)
	})

	r.Get("/active-counter", func(w http.ResponseWriter, r *http.Request) {
		res, err := db.Query("select count(*) from todos where completed = false;")
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		defer res.Close()

		var count int64
		res.Next()
		if err := res.Scan(&count); err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "active-counter.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "active-counter.tmpl", count)
	})

	r.Get("/clear-completed", func(w http.ResponseWriter, r *http.Request) {
		res, err := db.Query("select count(*) from todos where completed = true;")
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		defer res.Close()

		var count int64
		res.Next()
		if err := res.Scan(&count); err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "clear-completed.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "clear-completed.tmpl", count > 0)
	})

	r.Patch("/clear-completed", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec(`delete from todos where completed = true;`)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", "clearedCompleted")
		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "clear-completed.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "clear-completed.tmpl", false)
	})

	r.Get("/complete-all", func(w http.ResponseWriter, r *http.Request) {
		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		todosLeft := 0
		for _, todo := range todos {
			if !todo.Completed {
				todosLeft += 1
			}
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "complete-all.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "complete-all.tmpl", todosLeft == 0)
	})

	r.Put("/complete-all", func(w http.ResponseWriter, r *http.Request) {
		todos, err := allTodos(db)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}

		todosLeft := 0
		for _, todo := range todos {
			if !todo.Completed {
				todosLeft += 1
			}
		}

		newStatus := todosLeft > 0

		_, err = db.Exec("update todos set completed=(?)", newStatus)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("HX-Trigger", "completedAll")
		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "complete-all.tmpl"),
		}...))
		tmpl.ExecuteTemplate(w, "complete-all.tmpl", newStatus)
	})

	r.Get("/todos/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		todo, err := findTodo(id, db)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "todo.tmpl"),
		}...))
		data := Todo{Id: todo.Id, Todo: todo.Todo, Completed: todo.Completed, Editing: true}
		tmpl.ExecuteTemplate(w, "todo.tmpl", data)
	})

	r.Patch("/todos/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		r.ParseForm()
		form := r.Form
		todo := form.Get("todo")

		if todo == "" {
			_, err = db.Exec("delete from todos where id = (?);", id)
			if err != nil {
				http.Error(w, "fail", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			w.Header().Set("HX-Trigger", "deletedTodo")
			tmpl := template.Must(template.New("").Parse(""))
			tmpl.Execute(w, nil)

		} else {
			rows, err := db.Query("update todos set todo=(?) where id = (?) returning id, todo, completed;", todo, id)
			if err != nil {
				http.Error(w, "fail", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			defer rows.Close()

			var id int64
			var todo string
			var completed bool
			rows.Next()
			if err := rows.Scan(&id, &todo, &completed); err != nil {
				fmt.Println(err)
				http.Error(w, "fail", http.StatusInternalServerError)
				return
			}

			tmpl := template.Must(template.ParseFiles([]string{
				path.Join(basePath, "src", "todo.tmpl"),
			}...))
			data := Todo{Id: id, Todo: todo, Completed: completed, Editing: false}
			tmpl.ExecuteTemplate(w, "todo.tmpl", data)
		}
	})

	r.Delete("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		_, err = db.Exec("delete from todos where id = (?);", id)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("HX-Trigger", "deletedTodo")
		tmpl := template.Must(template.New("").Parse(""))
		tmpl.Execute(w, nil)
	})

	r.Patch("/todos/{id}/toggle", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		todo, err := findTodo(id, db)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		_, err = db.Exec("update todos set completed=(?) where id = (?);", !todo.Completed, id)
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		w.Header().Set("HX-Trigger", "updatedTodo")
		tmpl := template.Must(template.ParseFiles([]string{
			path.Join(basePath, "src", "todo.tmpl"),
		}...))
		data := Todo{Id: todo.Id, Todo: todo.Todo, Completed: !todo.Completed, Editing: false}
		tmpl.ExecuteTemplate(w, "todo.tmpl", data)
	})

	r.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
		sqlStmt := `
          drop table if exists todos;
          create table todos (id integer not null primary key, todo text, completed boolean default false);
        `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	fmt.Println("Server running on", "http://"+address)
	http.ListenAndServe(address, r)
}

func allTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query(`select id, todo, completed from todos order by id asc;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var id int64
		var todo string
		var completed bool

		if err := rows.Scan(&id, &todo, &completed); err != nil {
			return nil, err
		}

		todos = append(todos, Todo{Id: id, Todo: todo, Completed: completed})
	}

	return todos, nil
}

func whereTodos(completed bool, db *sql.DB) ([]Todo, error) {
	rows, err := db.Query(`select id, todo, completed from todos where completed = (?) order by id asc;`, completed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var id int64
		var todo string
		var completed bool

		if err := rows.Scan(&id, &todo, &completed); err != nil {
			return nil, err
		}

		todos = append(todos, Todo{Id: id, Todo: todo, Completed: completed, Editing: false})
	}

	return todos, nil
}

func findTodo(id int64, db *sql.DB) (*Todo, error) {
	res, err := db.Query("select id, todo, completed from todos where id = (?);", id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var todo string
	var completed bool
	res.Next()
	if err := res.Scan(&id, &todo, &completed); err != nil {
		return nil, err
	}

	return &Todo{Id: id, Todo: todo, Completed: completed, Editing: false}, nil
}
