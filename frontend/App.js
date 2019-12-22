import React from 'react';
import { Switch, Route } from 'react-router-dom';

import  PostsIndex from './resources/posts/Index';
import  PostsShow from './resources/posts/Show';

function App() {
  return (
    <div className="App">
      <h1>Hello React JS</h1>

      <Switch>
        <Route exact path="/posts" component={PostsIndex}></Route>
        <Route exact path="/posts/:id" component={PostsShow}></Route>
      </Switch>
    </div>
  )
}

export default App;
