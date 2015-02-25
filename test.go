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
    "encoding/json"
    //"time"
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
    go goio.listenForSockets()

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

    //daemon of incoming msg from the front end socket.
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
    rooms map[string][]*GoSocket

    idCount int
    goSockets map[int]*GoSocket

    socketsQueue chan *GoSocket

    roomActionsQueue chan RoomAction   
    
}

func NewGoio() *Goio {
    g := new(Goio)
    g.idCount = 0
    g.goSockets = make(map[int]*GoSocket)
    g.rooms = make(map [string] []*GoSocket)

    g.socketsQueue = make(chan *GoSocket)
    g.roomActionsQueue = make(chan RoomAction)

    return g
}

func (g *Goio) send(roomKey string, message string) {
    if slgs, ok := g.rooms[roomKey]; ok {
        for _, gsItem := range slgs {
            gsItem.conn.WriteMessage(websocket.TextMessage, []byte (message))
        }
    }
}


func (g *Goio) listenForSockets() {
    fmt.Println("listening on sockets...")
    defer close(g.socketsQueue)
    defer fmt.Println("exiting the listen task on sockets")

    for {
        x := <- g.socketsQueue
        g.idCount +=1
        x.id = g.idCount
        g.goSockets[x.id] = x
        x.listenConn()  //starts listening for further data

        fmt.Println("a socket gets retrieved from channel. id : %d", x.id)
    }
}


func (g *Goio) listenForRoomActions() {
    fmt.Println("listening for room actions...")
    defer close(g.roomActionsQueue)
    defer fmt.Println("exiting the listen task on room actions")

    for {
        ra := <- g.roomActionsQueue

        if _, ok := g.rooms[ra.roomKey]; ok { //operate action on existing slice
            switch (ra.action) {

            case ROOM_ACTION_ADD:
                //check that it is not already in
                a := g.rooms[ra.roomKey]
                found := false
                for _, gsItem :=range a {
                    if(gsItem == ra.goSocket) {
                        found = true
                    }
                }
                if (!found) {
                    g.rooms[ra.roomKey] = append(g.rooms[ra.roomKey], ra.goSocket)
                }

            case ROOM_ACTION_REMOVE:
                a := g.rooms[ra.roomKey]
                for i, gsItem := range a {
                    if gsItem == ra.goSocket {
                        a[i], a[len(a)-1], a = a[len(a)-1], nil, a[:len(a)-1]
                        g.rooms[ra.roomKey] = a
                        break
                    }
                }

            default:
                //do nothing.
            }

        } else { //create the slice in the map with 1 elt
            fmt.Println("could not find key 'inexistant");
        }
    }
}

func (g *Goio) getSocketById(id int) *GoSocket {
    if gs, ok := g.goSockets[id]; ok {
        return gs;
    }
    return nil;
}


type GoSocket struct {
    //conn wrapper
    conn *websocket.Conn

    goio *Goio

    references map[string] string

    id int

}

//daemon listening to incoming connections
func (gs *GoSocket) listenConn() {
    for {
        fmt.Println("reading...%d", count)
        messageType, msg, err := gs.conn.ReadMessage()
        if err != nil {
            fmt.Println("Failed to read a message: ", err)
            //NOT BREAK !! if BREAK, we should remove the socket from our rooms and watch.
            break
        }
        t := &Transport{}
        json.Unmarshal([]byte(str), &t)
    }
}

//to be used for such refs as user_id, etc...
func (gs *GoSocket) addReference(key string, value string) {
    gs.references[key] = value
}

func (gs *GoSocket) getReference(key string) string{
    return gs.references[key]
}

func (gs *GoSocket) isInRoom(roomKey string) bool {
    g := gs.goio
    if slgs, ok := g.rooms[roomKey]; ok {
        for _, gsItem := range slgs {
            if (gsItem == gs) {
                return true;
            }
        }
    }
    return false;
}

func (gs *GoSocket) broadcastSend(roomKey string, message string) {
    if gs.isInRoom(roomKey) {
        g := gs.goio
        if slgs, ok := g.rooms[roomKey]; ok {
            for _, gsItem := range slgs {
                if (gsItem != gs) {
                    gsItem.conn.WriteMessage(websocket.TextMessage, []byte (message))
                }
            }
        }
    }
}

//this is a concurrent method: many socket may join/leave at the exact same time; the join of 
//a socket should not block other joins/leave nor other actions (broadcast, emit.)
func (gs *GoSocket) join(roomKey string) {
    ra := RoomAction{roomKey, gs, ROOM_ACTION_ADD}
    gs.goio.roomActionsQueue <- ra
}

func (gs *GoSocket) leave(roomKey string) {
    ra := RoomAction{roomKey, gs, ROOM_ACTION_REMOVE}
    gs.goio.roomActionsQueue <- ra
}

const ROOM_ACTION_ADD int = 1;
const ROOM_ACTION_REMOVE int = -1;

type RoomAction struct {
    roomKey string
    goSocket *GoSocket
    action int
}

func NewGoSocket(conn *websocket.Conn, goio *Goio) *GoSocket {
    gs := new(GoSocket)
    gs.goio = goio
    gs.references = make(map [string] string)
    gs.conn = conn
    gs.goio.socketsQueue <- gs
    return gs
}

type Transport struct {
    EventName string        `json:"event_name"`
    Payload string          `json:"payload"`
    Conversation int        `json:"conversation"`
}

/*

func findStringInSlice(str string, sl *[]string) int{
    index := -1
    for i,element := range *sl {
        if(element == str) {
            return i;
        }
    }
    return index;
}



*/

