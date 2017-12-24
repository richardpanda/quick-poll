import React, { Component } from 'react';
import Button from 'react-toolbox/lib/button/Button';
import Card from 'react-toolbox/lib/card/Card';
import CardActions from 'react-toolbox/lib/card/CardActions';
import CardText from 'react-toolbox/lib/card/CardText';
import Checkbox from 'react-toolbox/lib/checkbox/Checkbox';
import Input from 'react-toolbox/lib/input/Input';

import './style.css';

import FormTitle from '../FormTitle';
import Loading from '../Loading';

import { fetchPostPoll } from '../../api';

class PollForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      checkIP: false,
      choices: ['', ''],
      error: '',
      isLoading: false,
      question: '',
    };
    this.handleChoiceChange = this.handleChoiceChange.bind(this);
    this.handleQuestionChange = this.handleQuestionChange.bind(this);
    this.handleCheckIPClick = this.handleCheckIPClick.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChoiceChange(index) {
    const self = this;
    return function(value) {
      const { choices } = self.state;
      const nextChoices = [...choices];
      const isLastChoice = index === choices.length - 1;
      const isChoiceEmpty = choices[index] === '';

      if (isLastChoice && isChoiceEmpty) {
        nextChoices.push('');
      }

      nextChoices[index] = value;
      self.setState({ choices: nextChoices });
    };
  }

  handleQuestionChange(value, { target }) {
    const { name } = target;
    this.setState({ [name]: value });
  }

  async handleSubmit(event) {
    event.preventDefault();

    const { history } = this.props;
    const { choices, checkIP, question } = this.state;
    const validChoices = choices.filter(choice => choice !== "");

    if (validChoices.length < 2) {
      this.setState({ error: 'Please provide at least two choices.' });
      return;
    }
    this.setState({ error: '', isLoading: true });

    try {
      const { id } = await fetchPostPoll({ question, choices: validChoices, checkIP });
      this.setState({ isLoading: false });
      history.push(`/polls/${id}`);
    } catch (e) {
      this.setState({ error: e.message, isLoading: false });
    }
  }

  handleCheckIPClick() {
    const { checkIP } = this.state;
    this.setState({ checkIP: !checkIP });
  }

  render() {
    const { checkIP, choices, error, isLoading, question } = this.state;

    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <Card style={{ margin: "16px auto", maxWidth: "600px" }}>
            <FormTitle
              title="Create a Poll"
              subtitle={error}
            />
            <CardText>
              <Input
                type="text"
                label="Question"
                name="question"
                value={question}
                onChange={this.handleQuestionChange}
                required
              />
              {choices.map((choice, i) => (
                <div key={i}>
                  <Input
                    type="text"
                    value={choice}
                    label="Choice"
                    onChange={this.handleChoiceChange(i)}
                  />
                </div>
              ))}
              <Checkbox
                className="poll-form__checkbox"
                label="IP Duplication Checking"
                checked={checkIP}
                onChange={this.handleCheckIPClick}
              />
            </CardText>
            <CardActions>
              {isLoading
                ? <Loading />
                : <Button type="submit" label="Submit" primary raised />
              }
            </CardActions>
          </Card>
        </form>
      </div>
    );
  }
}

export default PollForm;
