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
            if ( undefined !== ret.event_name && undefined !== ret.payload) {

                if (undefined !== ret.message_id) {
                    selfgoio.socket.defRecorder.resolveDef(ret.message_id, ret.event_name, ret.payload)
                }
                
                selfgoio.trigger(ret.event_name, ret.payload);
            } else {
                console.log("error in data format: could not find eventName and/or payload in server event");
            }
        })
    } else {
       /*WebSockets are not supported. Try a fallback method like long-polling etc*/
       this.trigger("error")
    }
}


Goio.prototype.emit=function(eventName, jsonData, optionalDef) {
    this.socket.send(eventName, jsonData, optionalDef);
}
////////////////////////////////////////////////////////


function DefRecorder() {
    this.defMap = {};
}

DefRecorder.prototype.defId = 0;

DefRecorder.prototype.registerDef = function(def) {
    console.log("registering def")
    DefRecorder.prototype.defId ++;
    var x = DefRecorder.prototype.defId;
    this.defMap[x] = def;
    return x;
}

DefRecorder.prototype.resolveDef = function(id, eventName, payload) {

    var def = this.defMap[id];
    if (undefined !== def && 0 !== id) {
        def.resolve(eventName, payload);
        delete this.defMap[id]; // pop out the reference.
    } else {
        console.log("callback for message_id '" + id + '"not found');
    }
}

//this is the wrapper for the websocket. it holds the websocket, and params for reconnection.
function GoioSocket(url, params) {
    this.websocket = null;
    //this.dispatcher = dispatcher;
    //todo: implement these parameters
    this.url = url;
    this['reconnection delay'] = 200;//random
    this['reconnection limit'] = 8000;
    this['max reconnection attempts'] = Infinity;
    this['auto connect'] = false;

    this.connected = false;
    this.defRecorder = new DefRecorder();

    _.extend(this, this.dispatcher);

    if (undefined === params) {
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


GoioSocket.prototype.send=function(eventName, jsonDataString, optionalDef) {

    if ('string' != typeof jsonDataString || 'string' != typeof jsonDataString) {
        console.log('event emit refused: type incorrect');
        return;
    }

    var messageId = 0; //0 stands for no message_id 
    if (undefined != optionalDef) {
        messageId = this.defRecorder.registerDef(optionalDef);
    }
    

    var pack = window.JSON.stringify({"event_name": eventName, "message_id": messageId, "payload":jsonDataString});

    console.log('sending...' + pack);
    this.websocket.send(pack);
    //error handling here

}




