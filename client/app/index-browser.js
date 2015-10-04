import React, {render} from 'react';
import createBrowserHistory from 'history/lib/createBrowserHistory';

import makeRoutes from './routes';

render(makeRoutes(createBrowserHistory()), document.getElementById("root"));
