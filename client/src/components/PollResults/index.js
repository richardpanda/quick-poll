import React, { Component } from 'react';
import Card from 'react-toolbox/lib/card/Card';
import CardText from 'react-toolbox/lib/card/CardText';
import CardTitle from 'react-toolbox/lib/card/CardTitle';

import './style.css';

import Loading from '../Loading';

class PollResults extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: {
        byId: {},
        allIds: [],
      },
      error: '',
      isLoading: false,
      question: '',
      ws: null,
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
        const byId = choices.reduce((acc, { id, text, num_votes }) => (
          { ...acc, [id]: { id, text, num_votes }}
        ), {});
        const allIds = choices.map(choice => choice.id);

        const ws = new WebSocket(`ws://localhost:8080/v1/polls/${id}/ws`);
        ws.onmessage = ({ data }) => {
          const { id, num_votes } = JSON.parse(data);
          const { choices } = this.state;
          if (choices.byId[id].num_votes < num_votes) {
            choices.byId[id].num_votes = num_votes;
            this.setState({ choices });
          }
        };

        this.setState({
          choices: { byId, allIds },
          isLoading: false,
          question,
          ws,
        });
      } else {
        this.setState({ error: payload.message, isloading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  componentWillUnmount() {
    this.state.ws.close();
  }

  render() {
    const { choices, error, isLoading, question } = this.state;
    const sum = choices.allIds.reduce((acc, id) => acc + choices.byId[id].num_votes, 0);
    const sortedChoices = choices.allIds.map(id => choices.byId[id]);
    sortedChoices.sort((c1, c2) => c2.num_votes - c1.num_votes);

    if (isLoading) {
      return <Loading className="loading-center" />
    }

    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <Card className="poll-results">
            <CardTitle title={question} />
            {sortedChoices.map(({ id, num_votes, text }) => (
              <CardText key={id}>
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
              </CardText>
            ))}
          </Card>
        </form>
      </div>
    );
  }
}

export default PollResults;
