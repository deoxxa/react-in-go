export const GET_CAT_PENDING = 'GET_CAT_PENDING';
export const GET_CAT_COMPLETE = 'GET_CAT_COMPLETE';
export const GET_CAT_FAILED = 'GET_CAT_FAILED';

export function getCat({ name }) {
  return (dispatch, getState) => {
    dispatch({
      type: GET_CAT_PENDING,
      meta: { name },
    });

    const { url } = getState().config.api;

    return fetch(`${url}/cats/${name}`).then((res) => res.json()).then((data) => {
      dispatch({
        GET_CAT_COMPLETE,
        meta: { name },
        data: data,
      });
    }).catch((err) => {
      dispatch({
        type: GET_CAT_FAILED,
        error: err,
      });
    });
  };
}

export const SEARCH_CATS_PENDING = 'SEARCH_CATS_PENDING';
export const SEARCH_CATS_COMPLETE = 'SEARCH_CATS_COMPLETE';
export const SEARCH_CATS_FAILED = 'SEARCH_CATS_FAILED';

export function searchCats({ search }) {
  return (dispatch, getState) => {
    dispatch({
      type: SEARCH_CATS_PENDING,
      meta: { search },
    });

    const { url } = getState().config.api;

    return fetch(`${url}/cats?search=${search}`).then((res) => res.json()).then((data) => {
      dispatch({
        SEARCH_CATS_COMPLETE,
        meta: { search },
        data: data,
      });
    }).catch((err) => {
      dispatch({
        type: SEARCH_CATS_FAILED,
        error: err,
      });
    });
  };
}
