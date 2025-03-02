package main

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

func main() {
    // Инициализируем gin.Engine
    r := gin.Default()

    // Загружаем шаблоны
    templateFiles, err := filepath.Glob("templates/*")
    if err != nil {
        log.Fatal(err)
    }
    r.SetHTMLTemplate(template.Must(template.ParseFiles(templateFiles...)))

    // Маршрутизация
    r.GET("/", indexHandler)
    r.StaticFile("/favicon.ico", "./static/favicon.ico")
    r.Static("/static", "./static")

    // Запуск сервера
    port := ":8080"
    log.Printf("Запуск сервера на порту %s\n", port)
    if err := r.Run(port); err != nil {
        log.Fatal(err)
    }
}

func indexHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "index.html", gin.H{})
}