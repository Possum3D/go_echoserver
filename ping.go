// server minimaliste
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {

    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.String(200, "We got Gin")
    })

    r.Run("127.0.0.1:8080")
}

