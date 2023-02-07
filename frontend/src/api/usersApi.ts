import { apiSlice } from './apiSlice';
import { User } from '../models/authUser';

const usersApi = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    getUsers: build.query<User[], void>({
      transformResponse(response: User[]) {
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
      providesTags: ['User'],
    }),
    getUser: build.query<User, string>({
      query(id: string) {
        return {
          url: `/users/${id}`,
          method: 'GET',
        };
      },
      providesTags: ['User'],
    }),
    postUser: build.mutation<User, Omit<User, 'id'>>({
      query: (user) => {
        return {
          url: '/user',
          method: 'POST',
          body: user,
        };
      },
      invalidatesTags: ['User'],
    }),
    putUser: build.mutation<User, User>({
      query: (user) => {
        return {
          url: `/users/${user.id}`,
          method: 'PUT',
          body: user,
        };
      },
      invalidatesTags: ['User'],
    }),
  }),
  overrideExisting: false,
});

export const { useGetUsersQuery, useGetUserQuery, usePostUserMutation, usePutUserMutation } = usersApi;
