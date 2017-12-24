import React, { Component } from 'react';
import Input from 'react-toolbox/lib/input/Input';
import Tooltip from 'react-toolbox/lib/tooltip';

const TooltipInput = Tooltip(Input);

class Share extends Component {
  constructor(props) {
    super(props);
    this.state = { tooltip: 'Copy to Clipboard' };
    this.handleClick = this.handleClick.bind(this);
  }

  componentDidMount() {
    this.input = document.getElementById("share-poll__input");
  }

  handleClick() {
    this.input.select();
    document.execCommand("Copy");
    this.setState({ tooltip: 'Copied!' });
  }

  render() {
    const { tooltip } = this.state;

    return (
      <div>
        <TooltipInput
          autoComplete="off"
          id="share-poll__input"
          label="Share"
          onClick={this.handleClick}
          onMouseLeave={() => setTimeout(() => this.setState({ tooltip: 'Copy to Clipboard' }), 200)}
          style={{ paddingBottom: "0px" }}
          tooltip={tooltip}
          tooltipHideOnClick={false}
          tooltipPosition="top"
          value={window.location.href}
        />
      </div>
    );
  }
}

export default Share;
