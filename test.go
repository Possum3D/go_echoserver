package main

//https://gowalker.org/github.com/gorilla/websocket
//http://golangtutorials.blogspot.fr/2011/06/web-programming-with-go-first-web-hello.html
//https://elithrar.github.io/article/custom-handlers-avoiding-globals/

import (
    "fmt"
    //"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "time"
)


type MyHttpHandler struct{}

func (h MyHttpHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request) {
    fmt.Fprint(w, "Hello!")
}



func main() {
    /*
    r := gin.Default()
    // this part is interesting if you want to enforce same origin policy./
    r.LoadHTMLFiles("goio.js") //load file
    r.LoadHTMLFiles("test.html") //load file


    r.GET("/", func(c *gin.Context) {
        c.HTML(200, "test.html", nil)
    })

    r.GET("/goio.js", func(c *gin.Context) {
        c.HTML(200, "goio.js", nil)
    })
    
    r.GET("/ws", func(c *gin.Context) {
        fmt.Println("get /ws called")
        wshandler(c.Writer, c.Request)
    })

    r.Run("0.0.0.0:8080")
    */


    //alternative à gin gonic.
    goio := NewGoio()
    go goio.listen()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        wshandler(w, r, goio)
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if ("" == r.URL.Path[1:]) {
            http.ServeFile(w, r, "test.html")
        } else {
            http.ServeFile(w, r, r.URL.Path[1:])
        }
        
    })

    err := http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        log.Fatal(err)
    }

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

func handler(w http.ResponseWriter, r *http.Request, mystr string) {
       println(mystr);
}

func wshandler(w http.ResponseWriter, r *http.Request, goio *Goio) {
    

    fmt.Println("upgrading a connection")
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: ", err)
        return
    }

    count :=0
    NewGoSocket(conn, goio)

    //boucle infinie = base du daemon server
    for {
        fmt.Println("reading...%d", count)
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Failed to read a message: ", err)
            break
        }

        //test input msg = Read with conn.ReadMessage
        msgstring := string(msg)

        if msgstring == "ping" {
            msgret := []byte("pong")
            conn.WriteMessage(messageType, msgret)
        } else {
            
            msgret := msg
            conn.WriteMessage(messageType, msgret)
        }
        count+=1
    }
}



//////////////////////////////////////////////////  goio    ////////////////////////////////////////////////////////////////
//serverside, GOIO is the parent that stores data about sockets
type Goio struct {
    //sockets will be added via a goroutine after it gets authentified.

    //ALL GoSockets are listed here, returned via sockets() method.
    //goSockets GoSockets
    //rooms: only a recorder of the rooms available, with the nb of GS associated.
    rooms map[string]int

    idCount int
    goSockets map[int]*GoSocket

    waitingSockets chan *GoSocket
    
    
}

func NewGoio() *Goio {
    g := new(Goio)
    g.idCount = 0
    g.goSockets = make(map[int]*GoSocket)

    g.waitingSockets = make(chan *GoSocket)

    return g
}

func (g *Goio) listen() {
    fmt.Println("listening on sockets...")
    defer close(g.waitingSockets)
    defer fmt.Println("listen is returning")

    for {
        x := <- g.waitingSockets
        g.idCount +=1
        x.id = g.idCount
        g.goSockets[x.id] = x
        fmt.Println("a socket gets retrieved from channel. id : %d", x.id)
    }
}


type GoSocket struct {
    //conn wrapper
    socket *websocket.Conn

    goio *Goio

    //slice aray of all rooms where the GS is.
    rooms []string

    id int

    //last action timestamp
    lastUpdated time.Time
}

func NewGoSocket(conn *websocket.Conn, goio *Goio) *GoSocket {
    gs := new(GoSocket)
    gs.goio = goio
    gs.socket = conn
    gs.goio.waitingSockets <- gs 
    return gs

}
/*
func (g *Goio) sockets() GoSockets {
    //return all sockets stored y Goio
    return g.goSockets;
}

func (g *Goio) addSocket(goSocket GoSocket) {
    g.goSockets.socket = append(g.goSockets.socket, goSocket)
}

func (g *Goio) removeSocket(goSocket GoSocket) {

}


//to be used like: goio.to('some room').emit('some event'):

func (g *Goio) to(goRoomKey string) GoSockets {
    x:= g.rooms[goRoomKey]
    

    sliceSock = []GoSocket
    s := GoSockets{sliceSock}

    if ( 0 == x) {
        return s;
    }
    
    for key, sock := range g.goSockets.sockets {
        if (findStringInSlice(goRoomKey, sock.rooms) != -1) {
            s.sockets = append(s.sockets, sock)
        }
    }
    return s;
}




func findStringInSlice(str string, sl *[]string) int{
    index := -1
    for i,element := range *sl {
        if(element == str) {
            return i;
        }
    }
    return index;
}

//to be used like: socket.join("some room")
func (gs *GoSocket) join(goRoomKey string) {
    //create if does not exist
    x := gs.goio.rooms[goRoomKey]

    if (x == 0) {
        gs.goio.rooms[goRoomKey] = 1


    } else {
        gs.goio.rooms[goRoomKey] +=1
    }
    gs.rooms = append(gs.rooms, goRoomKey)
    gs.goio.
}


func (gs *GoSocket) leave(goRoomKey string) {
    x := gs.goio.rooms[goRoomKey]
    if (x==0) {
        return;
    }

    //remove this key from slice.
    index := findStringInSlice(goRoomKey, gs.rooms)
    if (index == -1) {
        return;
    }
    gs.goio.rooms = append()

    //decrement count
    x -=1
}

type GoRoom struct {
    id string
    sockets GoSockets

}

//this is a bunch of sockets on which we can use emit
type GoSockets struct {
    sockets []GoSocket

}


//to be used like: goio.to("some room").emit("eventName", payload)
func (gss *GoSockets) emit() {
    //emit on every socket in the GoSockets bunch
}


//to be used like: goio.sockets.in("some room")
func (gss *GoSockets) in(goRoomKey string) GoSockets{

}

//returns GoSockets
func (gs * GoSocket) broadcastTo(goRoomKey string) GoSockets{
    //test that gs is part of the room

    //emit to all sockets except gs

}

*/

