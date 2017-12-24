import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Button from 'react-toolbox/lib/button/Button';
import Card from 'react-toolbox/lib/card/Card';
import CardActions from 'react-toolbox/lib/card/CardActions';
import CardText from 'react-toolbox/lib/card/CardText';
import RadioButton from 'react-toolbox/lib/radio/RadioButton';
import RadioGroup from 'react-toolbox/lib/radio/RadioGroup';

import './style.css';

import ErrorCard from '../ErrorCard';
import FormTitle from '../FormTitle';
import Loading from '../Loading';
import Share from '../Share';

class VoteForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: [],
      poll: {
        error: '',
        isFetching: false,
      },
      question: '',
      vote: {
        error: '',
        index: -1,
        isFetching: false,
      },
    };
    this.handleRadioChange = this.handleRadioChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  async componentDidMount() {
    const { id } = this.props.match.params;
    const endpoint = `/v1/polls/${id}`;
    this.setState({ poll: { isFetching: true, error: '' }});

    try {
      const response = await fetch(endpoint);

      if (response.status === 500) {
        this.setState({ poll: { error: 'Something went wrong...', isFetching: false }});
        return;
      }

      const payload = await response.json();

      if (response.ok) {
        const { choices, question } = payload;
        this.setState({ choices, poll: { error: '', isFetching: false }, question });
      } else {
        this.setState({ poll: { error: payload.message, isFetching: false }});
      }
    } catch (e) {
      this.setState({ poll: { error: e, isFetching: false }});
    }
  }

  handleRadioChange(index) {
    this.setState({ vote: { ...this.state.vote, index } });
  }

  async handleSubmit(event) {
    event.preventDefault();
    
    const { history, match } = this.props;
    const { choices, vote } = this.state;
    const { id: pollId } = match.params;
    if (vote.index === -1) {
      this.setState({ vote: { ...this.state.vote, error: 'Please select a choice.' }});
      return;
    }
    this.setState({ vote: { ...this.state.vote, error: '', isFetching: true }});

    try {
      const opts = { method: 'POST' };
      const { id: choiceId } = choices[vote.index];
      const response = await fetch(`/v1/polls/${pollId}/choices/${choiceId}`, opts);

      if (response.status === 500) {
        this.setState({ vote: { ...this.state.vote, error: 'Something went wrong...' }});
        return
      }

      const payload = await response.json();

      if (response.ok) {
        history.push(`/polls/${pollId}/results`);
      } else {
        this.setState({ vote: { ...this.state.vote, error: payload.message, isFetching: false }});
      }
    } catch (e) {
      this.setState({ vote: { ...this.state.vote, error: e, isFetching: false }});
    }
  }

  render() {
    const { match } = this.props;
    const { id } = match.params;
    const { choices, poll, question, vote } = this.state;

    if (poll.isFetching) {
      return (
        <Loading center />
      );
    }

    if (poll.error) {
      return <ErrorCard message="Something went wrong..." />;
    }

    return (
      <form onSubmit={this.handleSubmit}>
        <Card className="vote-form">
          <FormTitle title={question} subtitle={vote.error} />
          <CardText className="vote-form-radio-group">
            <RadioGroup name="vote" value={vote.index.toString()} onChange={this.handleRadioChange}>
              {choices.map((choice, i) => (
                <RadioButton key={choice.id} label={choice.text} value={i.toString()} />
              ))}
            </RadioGroup>
          </CardText>
          <CardActions className="vote-form-submit-button">
            {vote.isFetching
              ? <Loading />
              : (
                <div>
                  <Button type="submit" label="Vote" primary raised />
                  <Link to={`/polls/${id}/results`}>
                    <Button label="Results" accent raised />
                  </Link>
                </div>
              )}
              <Share />
          </CardActions>
        </Card>
      </form>
    );
  }
}

export default VoteForm;
