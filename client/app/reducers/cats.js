import { handleActions } from 'redux-actions';

import {
  GET_CAT_COMPLETE,
  SEARCH_CATS_COMPLETE,
} from '../actions';

export default handleActions({
  [GET_CAT_COMPLETE]: (state, action) => {
    return [action.data];
  },
  [SEARCH_CATS_COMPLETE]: (state, action) => {
    return action.data;
  },
}, []);
