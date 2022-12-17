import { apiSlice } from './apiSlice';
import { User } from '../models/authUser';

const usersApi = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    getUsers: build.query<User[], void>({
      transformResponse(response: User[], meta) {
        return response;
      },
      transformErrorResponse(response, meta) {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        const status = meta.response.status;
        return { response, status };
      },
      query() {
        return { url: '/users', method: 'GET' };
      },
    }),
  }),
  overrideExisting: false,
});

export const { useGetUsersQuery } = usersApi;
