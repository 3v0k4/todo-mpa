# TodoMPA

This repository is inspired by [TodoMVC](https://todomvc.com/), but comes with a twist: you won't find even one line of JavaScript here.

What's wild is that TodoMVC was created to highlight how frontend frameworks enable interactivity. However, this repo shows there's no need to pick such a framework to achieve the same level of app-iness™️.

TodoMPA is written in Go and, most importantly, [HTMX](https://htmx.org/).

Coded with ❤️  and 🤯 by [Riccardo Odone](https://odone.me).

## Additional Notes

You'll find two branches in this repository, where I explored the following question in HTMX:

> I need to update other content on the screen. How do I do this?

- `main` explores [*Solution 3: Triggering Events*](https://htmx.org/examples/update-other-content/#events)
- `oob` explores [*Solution 2: Out of Band Responses*](https://htmx.org/examples/update-other-content/#oob)

## Run

```bash
# install go
BASE_PATH="$(pwd)" go run src/main.go

# create database
curl -X POST http://127.0.0.1:3000/reset
```
