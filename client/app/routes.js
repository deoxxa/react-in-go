import React from 'react';
import {Router, Route, IndexRoute} from 'react-router';
import createMemoryHistory from 'history/lib/createMemoryHistory';

import {Layout, Search, Details} from './views';

export default (history) => <Router history={history || createMemoryHistory()}>
  <Route path="/" component={Layout}>
  	<IndexRoute component={Search} />
  	<Route path="/post/:id" component={Details} />
  </Route>
</Router>;
