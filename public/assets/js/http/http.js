'use default';

import {handleDebugErr} from '../errors/debug';

class HTTPClient {
    constructor(host, store){
        this.host = host;
        this.store = store;
    }
    request(method, path, body, extHdrs){
        let fullPath = `${this.host}${path}`;
        let req = new XMLHttpRequest();
        return new Promise((resolve, reject)=>{
            req.onload = this._onload.bind(this, req.responseText, resolve, reject);
            req.onerror = this._onerror.bind(this, reject);
            req.open(method.toUpperCase(), fullPath);
            if(body){
                req.send(body);
            } else {
                req.send();
            }
        });
    }
    handleDebugErr(){
        if(window.DEBUG){
            console.log('from HTTPClient: ');
        }
        handleDebugErr.apply(this, arguments);
    }
    get(path, body, extHdrs){
        return this.request('GET', this.fmtQueryParams(path, body), null, extHdrs);
    }
    post(path, body, extHdrs){
        return this.request('POST', path, body, extHdrs);
    }
    put(path, body, extHdrs){
        return this.request('PUT', path, body, extHdrs);
    }
    patch(path, body, extHdrs){
        return this.request('PATCH', path, body, extHdrs);
    }
    delete(path, body, extHdrs){
        return this.request('DELETE', this.fmtQueryParams(path, body), null, extHdrs);
    }
    options(path, extHdrs){
        return this.request('OPTIONS', path, null, extHdrs);
    }
    head(path, extHdrs){
        return this.request('HEAD', path, null, extHdrs);
    }
    fmtQueryParams(path, args){
        let keys = Object.keys(args);
        let pairs = keys.map((key)=>`${key}=${args[key]}`);
        if(pairs.length > 0){
            let res = `${path}?${pairs[0]}`;
            if(pairs.length > 1){
                pairs.slice(1).forEach((pair)=>{
                    res = `${res}&${pair}`;
                });
            }
            return res;
        }
        return path;
    }
    _onload(resp, resolve, reject){
        try{
            var msg = JSON.parse(resp);
        } catch(e) {
            console.error('HTTP response was not JSON');
            this.handleDebugErr(e);
            return reject(e);
        }
        return resolve(msg);
    }
    _onerror(reject, err){
        console.error('received an http error');
        this.handleDebugErr(err);
        return reject(err);
    }
}
