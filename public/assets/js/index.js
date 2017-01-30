'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, Link} from 'react-router';
import {createStore} from 'redux';
import Login from './login';

class Index extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            page: 'login'
        }
    }
    render(){
        return (
            <Router history={browserHistory}>
                <Route path="/" component={Login}/>
            </Router>
        );
    }
}

ReactDOM.render(<Index/>, documnet.getElemnetById('container'));
