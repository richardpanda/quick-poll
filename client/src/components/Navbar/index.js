import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import AppBar from 'react-toolbox/lib/app_bar/AppBar';

class Navbar extends Component {
  render() {
    return (
      <AppBar>
        <Link
          to="/"
          style={{
            fontSize: "18px",
            fontWeight: "700",
            textDecoration: "none",
          }}
        >
          Quick Poll
        </Link>
      </AppBar>
    );
  }
}

export default Navbar;
