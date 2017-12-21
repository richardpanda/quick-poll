import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import Button from 'react-toolbox/lib/button/Button';
import Card from 'react-toolbox/lib/card/Card';
import CardActions from 'react-toolbox/lib/card/CardActions';
import CardText from 'react-toolbox/lib/card/CardText';
import CardTitle from 'react-toolbox/lib/card/CardTitle';
import RadioButton from 'react-toolbox/lib/radio/RadioButton';
import RadioGroup from 'react-toolbox/lib/radio/RadioGroup';

import './style.css';

import Loading from '../Loading';

class VoteForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      choices: [],
      error: '',
      isLoading: true,
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
        this.setState({ choices, isLoading: false, question });
      } else {
        this.setState({ error: payload.message, isloading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  handleRadioChange(index) {
    this.setState({ vote: index });
  }

  async handleSubmit(event) {
    event.preventDefault();
    
    const { history, match } = this.props;
    const { choices, vote } = this.state;
    const { id: pollId } = match.params;
    if (vote === -1) {
      this.setState({ error: 'Please select a choice.' });
      return;
    }
    this.setState({ error: '', isLoading: true });

    try {
      const opts = { method: 'POST' };
      const { id: choiceId } = choices[vote];
      const response = await fetch(`/v1/polls/${pollId}/choices/${choiceId}`, opts);
      const payload = await response.json();

      if (response.ok) {
        this.setState({ isLoading: false });
        history.push(`/polls/${pollId}/results`);
      } else {
        this.setState({ error: payload.message, isLoading: false });
      }
    } catch (e) {
      this.setState({ error: e, isLoading: false });
    }
  }

  render() {
    const { match } = this.props;
    const { id } = match.params;
    const { choices, error, isLoading, question, vote } = this.state;

    if (isLoading) {
      return (
        <Loading className="loading-center" />
      );
    }

    return (
      <form onSubmit={this.handleSubmit}>
        <Card className="vote-form">
          <CardTitle title={question} subtitle={error} />
          <CardText className="vote-form-radio-group">
            <RadioGroup name="vote" value={vote.toString()} onChange={this.handleRadioChange}>
              {choices.map((choice, i) => (
                <RadioButton key={choice.id} label={choice.text} value={i.toString()} />
              ))}
            </RadioGroup>
          </CardText>
          <CardActions className="vote-form-submit-button">
            <Button type="submit" label="Vote" primary raised />
            <Link to={`/polls/${id}/results`}>
              <Button label="Results" accent raised />
            </Link>
          </CardActions>
        </Card>
      </form>
    );
  }
}

export default VoteForm;
