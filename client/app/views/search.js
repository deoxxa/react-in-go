import React, {Component} from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { searchCats } from '../actions';

@connect((state) => ({ cats: state.cats }), { searchCats })
export default class Search extends Component {
  static displayName = 'Search';

  static beforeRender(dispatch, params, query) {
    console.log('Search::beforeRender');

    return dispatch(searchCats({
      search: query.search || "",
    }));
  }

  componentDidMount() {
    this.props.searchCats({
      search: this.props.location.query.search || "",
    });
  }

  render() {
    return <pre>{this.props.location.pathname}</pre>;
  }
}
