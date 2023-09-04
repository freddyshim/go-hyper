# Go Hyper!
A "blazingly fast" full-stack web application in Go. This stack utilizes hypertext-driven templating tools such as [htmx](https://htmx.org) and Go's [html/template](https://pkg.go.dev/html/template) library to create highly interactive user interfaces without a single line of custom Javascript.

## The Stack
- [echo](https://github.com/labstack/echo) is the backbone, serving HTML templates via its' REST API.
- [htmx](https://htmx.org) allows HTML templates we send to the client to make requests to the server and populate sections of the client's page with interactive content. 
- [GORM](https://github.com/go-gorm/gorm) manages our schemas and connection to our PostgreSQL database.
- [TailwindCSS](https://github.com/tailwindlabs/tailwindcss) styles our HTML templates.

## Run Locally
From the root directory, create a `.env` file in-place and fill with with the following environment variables:

```
DATABASE_URL="postgres://username:password@localhost:5432/dbname"
```

Then, run the following command to start your server:

`go run main.go`

[air](https://github.com/cosmtrek/air) is recommended to use for hot-reloading. If you have it to run globally, run the command:

`air`
