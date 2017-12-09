import React, { Component } from 'react';
import {
  Card,
  CardPrimary,
  CardSubtitle,
  CardSupportingText,
  CardTitle,
  LinearProgress,
} from 'rmwc';

import './style.css';

class PollResults extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: [],
      error: '',
      isLoading: false,
      question: '',
    };
  }

  async componentDidMount() {
    this.setState({ isLoading: true });

    const { id } = this.props.match.params;
    const endpoint = `/v1/polls/${id}`;

    try {
      const response = await fetch(endpoint);
      const payload = await response.json();

      if (response.ok) {
        const { choices, question } = payload;
        this.setState({ choices, isLoading: false, question });
      } else {
        this.setState({ error: payload.message, isloading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  render() {
    const { choices, error, isLoading, question } = this.state;
    const sum = choices.reduce((acc, choice) => acc + choice.num_votes, 0);
    const sortedChoices = [...choices];
    sortedChoices.sort((c1, c2) => c2.num_votes - c1.num_votes);

    if (isLoading) {
      return <LinearProgress className="loading" determinate={false} />;
    }

    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <Card className="poll-results">
            <CardPrimary>
              <CardTitle large>{question}</CardTitle>
              <CardSubtitle className="poll-results-error">
                {error}
              </CardSubtitle>
            </CardPrimary>
            {sortedChoices.map(({ id, num_votes, text }) => (
              <CardSupportingText key={id}>
                <div className="poll-results-row">
                  <div>{text}</div>
                  {num_votes === 1
                    ? <div>{num_votes} Vote</div>
                    : <div>{num_votes} Votes</div>
                  }
                </div>
                <div className="poll-results-row">
                  <div className="poll-results-bar" style={{ width: `${num_votes / sum * 100}%` }} />
                  <div className="poll-results-percent">{Math.round(num_votes / sum * 100)}%</div>
                </div>
              </CardSupportingText>
            ))}
          </Card>
        </form>
      </div>
    );
  }
}

export default PollResults;
