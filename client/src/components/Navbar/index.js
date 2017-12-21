import React, { Component } from 'react';
import AppBar from 'react-toolbox/lib/app_bar/AppBar';
import Button from 'react-toolbox/lib/button/Button';
import Link from 'react-toolbox/lib/link/Link';
import Navigation from 'react-toolbox/lib/navigation/Navigation';

import './style.css';

class Navbar extends Component {
  render() {
    return (
      <AppBar title="Quick Poll" />
    );
  }
}

export default Navbar;
