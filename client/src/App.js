import React, { Component } from 'react';
import { Route, Switch } from 'react-router-dom';

import Navbar from './components/Navbar';
import PollForm from './components/PollForm';
import VoteForm from './components/VoteForm';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Navbar />
        <main>
          <Switch>
            <Route exact path="/" component={PollForm} />
            <Route path="/polls/:id" component={VoteForm} />
          </Switch>
        </main>
      </div>
    );
  }
}

export default App;
