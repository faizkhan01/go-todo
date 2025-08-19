package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Todo represents a todo item.
type Todo struct {
    ID      int64  `json:"id"`
    Title   string `json:"title"`
    Done    bool   `json:"done"`
    Details string `json:"details,omitempty"`
}

var (
    todos   = map[int64]*Todo{}
    nextID  int64 = 1
    todoMux sync.Mutex
)

func main() {
    r := gin.Default()

    // routes
    r.GET("/todos", listTodos)
    r.GET("/todos/:id", getTodo)
    r.POST("/todos", createTodo)
    r.PUT("/todos/:id", updateTodo)
    r.DELETE("/todos/:id", deleteTodo)

    // run server
    r.Run(":8080") // http://localhost:8080
}

// Handlers

// listTodos returns all todos
func listTodos(c *gin.Context) {
    todoMux.Lock()
    defer todoMux.Unlock()

    result := make([]*Todo, 0, len(todos))
    for _, t := range todos {
        result = append(result, t)
    }
    c.JSON(http.StatusOK, result)
}

// getTodo returns a todo by ID
func getTodo(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    todoMux.Lock()
    defer todoMux.Unlock()

    if t, ok := todos[id]; ok {
        c.JSON(http.StatusOK, t)
        return
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
}

// createTodo creates a new todo
func createTodo(c *gin.Context) {
    var req struct {
        Title   string `json:"title" binding:"required"`
        Details string `json:"details"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todoMux.Lock()
    defer todoMux.Unlock()

    t := &Todo{
        ID:      nextID,
        Title:   req.Title,
        Details: req.Details,
        Done:    false,
    }
    todos[nextID] = t
    nextID++

    c.JSON(http.StatusCreated, t)
}

// updateTodo updates an existing todo
func updateTodo(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req struct {
        Title   *string `json:"title"`
        Details *string `json:"details"`
        Done    *bool   `json:"done"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todoMux.Lock()
    defer todoMux.Unlock()

    t, ok := todos[id]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
        return
    }

    if req.Title != nil {
        t.Title = *req.Title
    }
    if req.Details != nil {
        t.Details = *req.Details
    }
    if req.Done != nil {
        t.Done = *req.Done
    }

    c.JSON(http.StatusOK, t)
}

// deleteTodo removes a todo
func deleteTodo(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    todoMux.Lock()
    defer todoMux.Unlock()

    if _, ok := todos[id]; !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
        return
    }
    delete(todos, id)
    c.Status(http.StatusNoContent)
    c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}