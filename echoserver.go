package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "net/http"
)

func main() {

    r := gin.Default()
    // this part is interesting if you want to enforce same origin policy./
    r.LoadHTMLFiles("index.html") //load file

    r.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", nil)
    })
    
    r.GET("/ws", func(c *gin.Context) {
        wshandler(c.Writer, c.Request)
    })

    r.Run("0.0.0.0:8080")
}


//class implémentant Upgrader.
var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    //add the following if you do NOT want to check origin. otherwise:default checks same origin.
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: ", err)
        return
    }

    //boucle infinie = base du daemon server
    for {
        t, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Failed to read a message: ", err)
            break
        }

        //test input msg = Read with conn.ReadMessage
        msgstring := string(msg)

        if msgstring == "ping" {
          msgret := []byte("pong")
          conn.WriteMessage(t, msgret)
        } else {
          msgret := msg
          conn.WriteMessage(t, msgret)
        }
        
    }
}


/*package main

import (
  "fmt"
  "github.com/gorilla/websocket"
  "net/http"
  "github.com/gin-gonic/gin"
)

var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: %+v", err)
        return
    }

    for {
        t, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        conn.WriteMessage(t, msg)
    }
}

func main() {
  fmt.Println("Starting websock server: ")
  http.Handle("/echo", websocket.Handler(webHandler))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("ListenAndServe: " + err.Error())
    //log.Fatal("ListenAndServe:", err)
  }
}

/*

func main() {
  flag.Parse()
  homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
  go h.run()
  http.HandleFunc("/", homeHandler)
  http.HandleFunc("/ws", wsHandler)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
*/
