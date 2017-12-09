import React, { Component } from 'react';
import {
  Card,
  CardAction,
  CardActions,
  CardPrimary,
  CardSubtitle,
  CardSupportingText,
  CardTitle,
  Radio,
} from 'rmwc';

import './style.css';

class VoteForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: [],
      error: '',
      question: '',
      vote: -1,
    };
    this.handleRadioChange = this.handleRadioChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  async componentDidMount() {
    const { id } = this.props.match.params;
    const endpoint = `/v1/polls/${id}`;

    try {
      const response = await fetch(endpoint);
      const payload = await response.json();

      if (response.ok) {
        const { choices, question } = payload;
        this.setState({ choices, question });
      } else {
        this.setState({ error: payload.message });
      }
    } catch (e) {
      this.setState({ error: e });
    }
  }

  handleRadioChange(index) {
    const self = this;
    return function(event) {
      self.setState({ vote: index });
    };
  }

  async handleSubmit(event) {
    event.preventDefault();
    
    const { choices, vote } = this.state;
    if (vote === -1) {
      this.setState({ error: 'Please select a choice.' });
      return;
    }
    this.setState({ error: '' });

    try {
      const opts = { method: 'POST' };
      const id = choices[vote].id;
      const response = await fetch(`/v1/choices/${id}`, opts);
      const payload = await response.json();

      if (response.ok) {
        console.log(payload);
      } else {
        this.setState({ error: payload.message });
      }
    } catch (e) {
      this.setState({ error: e });
    }
  }

  render() {
    const { choices, error, question, vote } = this.state;

    return (
      <form onSubmit={this.handleSubmit}>
        <Card className="vote-form">
          <CardPrimary className="vote-form-title">
            <CardTitle large>{question}</CardTitle>
            <CardSubtitle className="vote-form-error">
              {error}
            </CardSubtitle>
          </CardPrimary>
          <CardSupportingText>
            {choices.map((choice, i) => (
              <div key={choice.id}>
                <Radio
                  name="radio"
                  checked={vote === i}
                  onChange={this.handleRadioChange(i)}
                >
                  {choice.text}
                </Radio>
              </div>
            ))}
          </CardSupportingText>
          <CardActions>
            <CardAction unelevated type="submit">Vote</CardAction>
          </CardActions>
        </Card>
      </form>
    );
  }
}

export default VoteForm;
