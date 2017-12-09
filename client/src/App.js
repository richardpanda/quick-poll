import React, { Component } from 'react';
import { Route, Switch } from 'react-router-dom';

import Navbar from './components/Navbar';
import PollForm from './components/PollForm';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Navbar />
        <main>
          <Switch>
            <Route exact path="/" component={PollForm} />
          </Switch>
        </main>
      </div>
    );
  }
}

export default App;
