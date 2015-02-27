package main

//https://gowalker.org/github.com/gorilla/websocket
//http://golangtutorials.blogspot.fr/2011/06/web-programming-with-go-first-web-hello.html
//https://elithrar.github.io/article/custom-handlers-avoiding-globals/

import (
    "fmt"
    //"github.com/gin-gonic/gin"
    "net/http"
    "log"
    "strconv"
    "github.com/mentionapp/common.go/goio"
    //"time"
)


type MyHttpHandler struct{}

func (h MyHttpHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request) {
    fmt.Fprint(w, "Hello!")
}


func main() {

    g := goio.NewGoio()

    defer g.ShutDown()
    log.Println("test de logprint")


    g.OnConnexion(func(gs *goio.GoSocket){
        fmt.Println("in onConnexion!!")
        //gs.join(strconv.Itoa(gs.id))
        gs.Join("testRoom"+strconv.Itoa(gs.Id))
        //fmt.Println("room joined: ", strconv.Itoa(gs.id))
        fmt.Println("room joined: ", "testRoom"+ strconv.Itoa(gs.Id))

        g.On("echo", func(ev goio.FrontEvent, v ...interface{}){
            if (ev.GoSocket == gs) {
                fmt.Println("callback called!! sendint ")
                //goio.send(strconv.Itoa(gs.id), ev)
                //g.Send("testRoom" + strconv.Itoa(gs.Id), ev)
                g.ReplyTo(ev, "{\"ping\":\"replied!\"}")
            }
            
        })
        fmt.Println("callback added")
    })

    g.HandleWebsocket("/ws")


    //for tests only..

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if ("" == r.URL.Path[1:]) {
            http.ServeFile(w, r, "test.html")
        } else {
            http.ServeFile(w, r, r.URL.Path[1:])
        }
    })

    //end for tests only.

    err := http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        log.Fatal(err)
    }

}



///////////////////////////// TESTS (to be done)
// what if sent json is invalid
// what if sent json does not include expected data
// what if a socket connexion is lost
// is it possible that response is linked to wrong question after connection was lost / re done
// is there a possibility of memory leak through closed connections not removed from rooms, etc...


///////////////////////////// TODO
// put goio part in separate package (as : the middleware)
// test goio package.



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

