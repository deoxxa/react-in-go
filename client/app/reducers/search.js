import { handleActions } from 'redux-actions';
import { mergeIn, setIn } from 'nonmutable';

import {
  GET_CAT_PENDING,
  GET_CAT_COMPLETE,
  GET_CAT_ERROR,
  SEARCH_CATS_PENDING,
  SEARCH_CATS_COMPLETE,
  SEARCH_CATS_ERROR,
} from '../actions';

const initialState = {
  loading: false,
  error: null,
};

export default handleActions({
  [GET_CAT_PENDING]: (state, action) => {
    return {
      loading: true,
      error: null,
    };
  },
  [GET_CAT_COMPLETE]: (state, action) => {
    return {
      loading: false,
      error: null,
    };
  },
  [GET_CAT_ERROR]: (state, action) => {
    return {
      loading: false,
      error: action.error,
    };
  },
  [SEARCH_CATS_PENDING]: (state, action) => {
    return {
      loading: true,
      error: null,
    };
  },
  [SEARCH_CATS_COMPLETE]: (state, action) => {
    return {
      loading: false,
      error: null,
    };
  },
  [SEARCH_CATS_ERROR]: (state, action) => {
    return {
      loading: false,
      error: action.error,
    };
  },
}, initialState);
