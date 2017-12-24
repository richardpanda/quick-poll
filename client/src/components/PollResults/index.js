import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Button from 'react-toolbox/lib/button/Button';
import Card from 'react-toolbox/lib/card/Card';
import CardActions from 'react-toolbox/lib/card/CardActions';
import CardText from 'react-toolbox/lib/card/CardText';
import CardTitle from 'react-toolbox/lib/card/CardTitle';

import './style.css';

import ErrorCard from '../ErrorCard';
import Loading from '../Loading';
import Share from '../Share';

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

      if (response.status === 500) {
        this.setState({ error: 'Something went wrong...', isLoading: false });
        return;
      }

      const payload = await response.json();

      if (response.ok) {
        const { choices, question } = payload;
        const byId = choices.reduce((acc, { id, text, num_votes }) => (
          { ...acc, [id]: { id, text, num_votes }}
        ), {});
        const allIds = choices.map(choice => choice.id);

        let ws = null;
        if (process.env.NODE_ENV === "production") {
          ws = new WebSocket(`ws://${window.location.hostname}/v1/polls/${id}/ws`);
        } else {
          ws = new WebSocket(`ws://localhost:8080/v1/polls/${id}/ws`);
        }

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
        this.setState({ error: payload.message, isLoading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  componentWillUnmount() {
    if (this.state.ws) {
      this.state.ws.close();
    }
  }

  render() {
    const { match } = this.props;
    const { choices, error, isLoading, question } = this.state;
    const sum = choices.allIds.reduce((acc, id) => acc + choices.byId[id].num_votes, 0);
    const sortedChoices = choices.allIds.map(id => choices.byId[id]);
    sortedChoices.sort((c1, c2) => c2.num_votes - c1.num_votes);

    if (isLoading) {
      return <Loading center />
    }

    if (error) {
      return <ErrorCard message={error} />; 
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
                  <div className="poll-results-percent">
                    {sum === 0 ? 0 : Math.round(num_votes / sum * 100)}%
                  </div>
                </div>
              </CardText>
            ))}
            <CardActions style={{ height: "48px", justifyContent: "space-between" }}>
              <div>
                <Link to={`/polls/${match.params.id}`}>
                  <Button label="Vote Page" primary raised style={{ marginRight: "4px" }} />
                </Link>
                <Link to={`/`}>
                  <Button label="New Poll" accent raised />
                </Link>
              </div>
              <Share />
            </CardActions>
          </Card>
        </form>
      </div>
    );
  }
}

export default PollResults;
