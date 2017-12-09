import React, { Component } from 'react';
import {
  Card,
  CardAction,
  CardActions,
  CardPrimary,
  CardSubtitle,
  CardSupportingText,
  CardTitle,
  LinearProgress,
  TextField,
} from 'rmwc';

import './style.css';

class PollForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: ['', ''],
      error: '',
      isLoading: false,
      question: '',
    };
    this.handleChoiceChange = this.handleChoiceChange.bind(this);
    this.handleQuestionChange = this.handleQuestionChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChoiceChange(index) {
    const self = this;
    return function({ target }) {
      const { value } = target;
      const { choices } = self.state;
      const nextChoices = [...choices];

      if (index === choices.length-1 && choices[index] === '') {
        nextChoices.push('');
      }

      nextChoices[index] = value;
      self.setState({ choices: nextChoices });
    };
  }

  handleQuestionChange({ target }) {
    const { value } = target;
    this.setState({ question: value });
  }

  async handleSubmit(event) {
    event.preventDefault();

    const { history } = this.props;
    const { choices, question } = this.state;
    const validChoices = choices.filter(choice => choice !== "");

    if (validChoices.length < 2) {
      this.setState({ error: 'Please provide at least two choices.' });
      return;
    }
    this.setState({ error: '', isLoading: true });

    try {
      const opts = {
        method: 'POST',
        body: JSON.stringify({ question, choices: validChoices }),
      };
      const response = await fetch('/v1/polls', opts);
      const payload = await response.json();

      if (response.ok) {
        this.setState({ isLoading: false });
        history.push(`/polls/${payload.id}`);
      } else {
        this.setState({ error: payload.message, isLoading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  render() {
    const { choices, error, isLoading } = this.state;

    return (
      <div>
        {isLoading
          && <LinearProgress className="loading" determinate={false} />
        }
        <form onSubmit={this.handleSubmit}>
          <Card className="poll-form">
            <CardPrimary>
              <CardTitle large>Create a Poll</CardTitle>
              <CardSubtitle className="poll-form-error">
                {error}
              </CardSubtitle>
            </CardPrimary>
            <CardSupportingText>
              <div>
                <TextField
                  className="text-field"
                  label="Question"
                  onChange={this.handleQuestionChange}
                  required
                />
              </div>
              {choices.map((choice, i) => (
                <div key={i}>
                  <TextField
                    className="text-field"
                    label="Choice"
                    onChange={this.handleChoiceChange(i)}
                  />
                </div>
              ))}
            </CardSupportingText>
            <CardActions>
              <CardAction unelevated type="submit">Submit</CardAction>
            </CardActions>
          </Card>
        </form>
      </div>
    );
  }
}

export default PollForm;
