'use strict';

import {handleDebugErr} from '../errors/debug'

class WSClient extends WebSocket {
    constructor(host, protocols, store){
        super(host, protocols);
        this._host = host;
        this.store = store;
        this.subscription = null;
    }
    handleDebugErr(){
        console.log('from WSClient: ');
        handleDebugErr.apply(this, arguments);
    }
    onopen(){
        if(window.DEBUG){
            console.log(`connection to ${this._host} is open`);
        }
        this.store.subscribe(()=>{
            if(this.subscription){
                this.unsubscribeTo(this.subscription);
            }
            let currState = this.store.getState();
            this.subscribeTo(currState.page);
        });
    }
    onmessage(message){
        try{
            var msg = JSON.parse(message);
        } catch(e) {
            console.error('websocket message is not JSON');
            this.handleDebugErr(e);
        }
        for(let i in msg){
            if((i == 'project' || i == 'branch') && msg[i] != this.subscription[i]){
                return;
            }
        }
        store.dispatch(msg);
    }
    onerror(e){
        console.error('an error has occurred on the websocket');
        this.handleDebugErr(e);
    }
    send(msg){
        super.send(JSON.stringify(msg));
    }
    subscribeTo(channel){
        let payload = Object.assign({}, channel, {
            action: 'subscribe',
            innerMessage: {
                timestamp: new Date().toUTCString()
            }
        });
        this.send(payload);
        this.subscription = channel;
    }
    unsubscribeTo(channel){
        let payload = Object.assign({}, channel, {
            action: 'unsubscribe',
            innerMessage: {
                timestamp: new Date().toUTCString()
            }
        });
        this.send(payload);
        this.subscription = null;
    }
    close(){
        this.unsubscribeTo(this.subscription);
        super.close();
    }
}
