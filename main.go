package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Todo struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Content string
}

func main() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Todo{})

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	tmpl := template.Must(
		template.Must(template.ParseGlob("view/*.html")).ParseGlob("view/icons/*.svg"))
	e.Renderer = &Template{
		templates: tmpl,
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/todo", func(c echo.Context) error {
		var todos []Todo
		result := db.Order("created_at asc").Find(&todos)
		if result.Error != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
		}
		return c.Render(http.StatusOK, "todos", todos)
	})

	e.POST("/todo", func(c echo.Context) error {
		todo := Todo{
			Content: "",
		}
		db.Create(&todo)
		return c.Render(http.StatusOK, "todo-edit", todo)
	})

	e.GET("/todo/:id", func(c echo.Context) error {
		edit := c.QueryParam("edit")
		id, error := uuid.Parse(c.Param("id"))
		if error != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid id"}
		}

		todo := Todo{
			ID: id,
		}
		result := db.First(&todo)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid id"}
		}

		if edit == "true" {
			return c.Render(http.StatusOK, "todo-edit", todo)
		}
		return c.Render(http.StatusOK, "todo", todo)
	})

	e.PUT("/todo/:id", func(c echo.Context) error {
		id, error := uuid.Parse(c.Param("id"))
		content := c.FormValue("content")
		if error != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid id"}
		}

		todo := Todo{
			ID: id,
		}
		result := db.Model(&todo).Update("content", content)
		if result.Error != nil {
			fmt.Println(result.Error)
		}

		return c.Render(http.StatusOK, "todo", todo)
	})

	e.DELETE("/todo/:id", func(c echo.Context) error {
		id, error := uuid.Parse(c.Param("id"))
		if error != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid id"}
		}

		todo := Todo{
			ID: id,
		}
		result := db.Delete(&todo)
		if result.Error != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "unable to delete todo"}
		}

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
