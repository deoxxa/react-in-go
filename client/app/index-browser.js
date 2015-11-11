import React, {render} from 'react';
import createBrowserHistory from 'history/lib/createBrowserHistory';
import { Provider } from 'react-redux';
import { merge } from 'nonmutable';

import makeRoutes from './routes';
import createStore from './store';

const store = createStore(merge(window.initialState || {}, {
  config: {
    api: {
      url: `${document.location.origin}/api/v1`,
      jwt: '',
    },
  },
}));

render(
	<Provider store={store}>
		{() => makeRoutes(createBrowserHistory())}
	</Provider>
, document.getElementById("root"));
