import React, { Component } from 'react';

import Navbar from './components/Navbar';
import PollForm from './components/PollForm';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Navbar />
        <PollForm />
      </div>
    );
  }
}

export default App;
