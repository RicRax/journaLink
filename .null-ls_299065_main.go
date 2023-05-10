p
ackage main

import(
    "net/http"
    "github.com/gin-gonic/gin"
)

type todo struct{
    ID          string `json:"id"`
    Item        string `json:"title"`
    Completed   bool   `json:"completed"`
}

var todos = []todo{
    {ID: "1", Item: "Clean room", Completed: false},
    {ID: "1", Item: "Clean room", Completed: false},
    {ID: "1", Item: "Clean room", Completed: false},
}

func main() {
    router:= gin.Default()
    router.GET("/todos",)
        router.Run("localhost:8080")
}
