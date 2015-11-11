import React, {Component} from 'react';
import { connect } from 'react-redux';
import { getCat } from '../actions';

@connect((state) => ({ cats: state.cats }), { getCat })
export default class Details extends Component {
  static displayName = 'Details';

  static beforeRender(dispatch, params) {
    console.log('Details::beforeRender');

    return dispatch(getCat({
      name: params.name,
    }));
  }

  componentDidMount() {
    this.props.getCat({
      name: this.props.params.name,
    });
  }

  componentWillReceiveProps(newProps) {
    if (newProps.params.name !== this.props.params.name) {
      this.props.getCat({
        name: this.props.params.name,
      });
    }
  }

  render() {
    return <pre>{this.props.location.pathname}</pre>;
  }
}
