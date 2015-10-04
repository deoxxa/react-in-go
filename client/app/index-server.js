import React, {renderToString} from 'react';
import { RoutingContext, match } from 'react-router';
import createLocation from 'history/lib/createLocation';

import makeRoutes from './routes';

const routes = makeRoutes();

renderApplication = (url, res) => {
  const location = createLocation(url);

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

    return res.success(renderToString(<RoutingContext {...renderProps} />));
  });
};
