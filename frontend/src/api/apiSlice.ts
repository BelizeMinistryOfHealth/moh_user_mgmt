import { BaseQueryFn, createApi, FetchArgs, fetchBaseQuery, FetchBaseQueryError } from '@reduxjs/toolkit/query/react';
import { STORAGE_KEYS } from '../appConstants';

const baseQuery = fetchBaseQuery({
  baseUrl: 'https://users-mgmt-e46d3zpgka-ue.a.run.app',
  prepareHeaders(headers) {
    const state = localStorage.getItem(STORAGE_KEYS.USER_DATA);
    if (state) {
      const userState: { user: { token: string } } = JSON.parse(state);
      headers.set('Authorization', `Bearer ${userState.user.token}`);
    }
    return headers;
  },
});

const baseQueryWithAuthCheck: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions,
) => {
  const result = await baseQuery(args, api, extraOptions);
  const status = result?.meta?.response?.status;
  if (status == 401) {
    localStorage.removeItem(STORAGE_KEYS.USER_DATA);
  }
  return result;
};
export const apiSlice = createApi({
  reducerPath: 'apiSlice',
  baseQuery: baseQueryWithAuthCheck,
  endpoints: () => ({}),
  tagTypes: ['User'],
});
