import { createStore, applyMiddleware, combineReducers } from 'redux';
import thunk from 'redux-thunk';
import * as reducers from './reducers';

const createStoreWithMiddleware = applyMiddleware(thunk)(createStore);

export default function(initialState = {}) {
	return createStoreWithMiddleware(combineReducers(reducers), initialState);
}
