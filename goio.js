//dependencies:
//backbone, underscore, jquery.
//credits :http://www.developerfusion.com/article/143158/an-introduction-to-websockets/

function goio() {
    
    var g = new Goio();

    return g;
}

function Goio (dispatcher) {
    this.socket = null;
}
_.extend(Goio.prototype, Backbone.Events);


Goio.prototype.connect=function(url, params) {
    if ('WebSocket' in window){
       /* WebSocket is supported. You can proceed with your code*/
       
       this.socket = new GoioSocket(url, params);
       
        this.trigger("connect");
        selfgoio = this;

        this.socket.on("message", function(msgEvent){

            var ret = JSON.parse(msgEvent.data);
            //console.log('ret: '.ret);
            if ( undefined !==ret.eventName && undefined !== ret.payload) {
                selfgoio.trigger(ret.eventName, ret.payload);
            } else {
                console.log("error in data format: could not find eventName and/or payload in server event");
            }
        })

    } else {
       /*WebSockets are not supported. Try a fallback method like long-polling etc*/
       this.trigger("error")
    }
}


Goio.prototype.emit=function(eventName, jsonData) {
    this.socket.send(eventName, jsonData);
}
////////////////////////////////////////////////////////


//this is the wrapper for the websocket. it holds the websocket, and params for reconnection.
function GoioSocket(url, params) {
    this.websocket = null;
    //this.dispatcher = dispatcher;
    //todo: implement these parameters
    this.url = url;
    this['reconnection delay'] = 200;
    this['reconnection limit'] = 8000;
    this['max reconnection attempts'] = Infinity;
    this['auto connect'] = false;

    this.connected = false;

    _.extend(this, this.dispatcher);

    if (undefined ===params) {
        return;
    }

    Object.keys(params).forEach(function(key,value) {
        if(this.hasOwnProperty(key)) {
            this[key] = value;
        }
    })

}

_.extend(GoioSocket.prototype, Backbone.Events);


GoioSocket.prototype.connect=function() {
    goioSocketSelf = this;
    this.websocket = new WebSocket(this.url);
    this.websocket.onopen = function(){
        this.connected = true;
        goioSocketSelf.trigger('connection');       
    }

    this.websocket.onclose = function(){
        this.connected = false;
        goioSocketSelf.trigger('closed');
    }

    this.websocket.onmessage = function(msgEvent){

        goioSocketSelf.trigger("message", msgEvent);
    }

    this.websocket.onerror = function(){
        goioSocketSelf.trigger('error');
    }
}


/*
GoioSocket.prototype.on=function(eventName, callback) {
    
    this.websocket.onmessage = function(serverMsg){
        //extract json : should have an eventName and a body, from serverMsg. mon event should be extracted from serverMsg.
        data = {"eventName": "monevent", "payload":{}};
        if(data.eventName == eventName) {
            callback.apply(arguments)
        }
    }

}*/

GoioSocket.prototype.send=function(eventName, jsonDataString) {

    var pack = window.JSON.stringify({"eventName": eventName, "payload":jsonDataString});

    console.log('sending...' + pack);
    this.websocket.send(pack);
    //error handling here

}




