import React, {Component} from 'react';

export default class Details extends Component {
  render() {
    return <div>{JSON.stringify(this.props.params)}</div>;
  }
}
