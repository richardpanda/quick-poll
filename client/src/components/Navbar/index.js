import React, { Component } from 'react';
import AppBar from 'react-toolbox/lib/app_bar/AppBar';

import './style.css';

class Navbar extends Component {
  render() {
    return (
      <AppBar title="Quick Poll" />
    );
  }
}

export default Navbar;
