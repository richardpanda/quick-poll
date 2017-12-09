import React, { Component } from 'react';
import { Toolbar, ToolbarRow, ToolbarTitle } from 'rmwc';

class Navbar extends Component {
  render() {
    return (
      <Toolbar>
        <ToolbarRow>
          <ToolbarTitle>Quick Poll</ToolbarTitle>
        </ToolbarRow>
      </Toolbar>
    );
  }
}

export default Navbar;
