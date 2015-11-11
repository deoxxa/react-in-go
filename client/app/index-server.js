import React, {renderToString} from 'react';
import { RoutingContext, match } from 'react-router';
import createLocation from 'history/lib/createLocation';
import { Provider } from 'react-redux';
import { parse } from 'qs';

import makeRoutes from './routes';
import createStore from './store';

const routes = makeRoutes();

renderApplication = (url, ctx, res) => {
  // apparently `query` isn't populated automatically - not sure what's up!
  const location = createLocation(url);
  location.query = location.search ? parse(location.search.substr(1)) : {};

  const store = createStore({
    config: {
      api: {
        url: '',
      },
    },
  });

  return match({routes, location}, (error, redirectLocation, renderProps) => {
    if (redirectLocation) {
      return res.redirect(redirectLocation.pathname + redirectLocation.search);
    }

    if (error) {
      return res.error(error.message);
    }

    if (renderProps == null) {
      return res.notFound();
    }

    const promises = renderProps.routes.map((route) => {
      if (!route.component) {
        return null;
      }

      if (typeof route.component.beforeRender !== 'function') {
        return null;
      }

      return route.component.beforeRender(::store.dispatch, renderProps.params, renderProps.location.query);
    }).filter((e) => !!e);

    Promise.all(promises).then(() => {
      return res.success(JSON.stringify(store.getState()), renderToString(
        <Provider store={store}>
          {() => <RoutingContext {...renderProps} />}
        </Provider>
      ));
    });
  });
};
