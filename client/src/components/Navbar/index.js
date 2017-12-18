import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Toolbar, ToolbarRow, ToolbarTitle } from 'rmwc';

import './style.css';

class Navbar extends Component {
  render() {
    return (
      <Toolbar>
        <ToolbarRow>
          <ToolbarTitle><Link className="home-link" to="/">Quick Poll</Link></ToolbarTitle>
        </ToolbarRow>
      </Toolbar>
    );
  }
}

export default Navbar;
