import { apiSlice } from './apiSlice';
import { RawUser, User } from '../models/authUser';
import { UserApplication } from '../models/userApplications';

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
      transformResponse(response: RawUser) {
        return {
          ...response,
          userApplications: response.userApplications?.find((app) => app.name === 'hiv_surveys') ?? null,
        };
      },
      query(id: string) {
        return {
          url: `/users/${id}`,
          method: 'GET',
        };
      },
      providesTags: ['User'],
    }),
    getApplications: build.query<UserApplication | null | undefined, void>({
      transformResponse: (response: UserApplication[]) => {
        return response.find((app) => app.name === 'hiv_surveys');
      },
      query() {
        return {
          url: '/applications',
          method: 'GET',
        };
      },
    }),
    postUser: build.mutation<User, Omit<RawUser, 'id'>>({
      transformResponse(response: RawUser) {
        return {
          ...response,
          userApplications: response.userApplications?.find((app) => app.name === 'hiv_surveys') ?? null,
        };
      },
      query: (user) => {
        return {
          url: '/user',
          method: 'POST',
          body: user,
        };
      },
      invalidatesTags: ['User'],
    }),
  }),
  overrideExisting: false,
});

export const { useGetUsersQuery, useGetApplicationsQuery, useGetUserQuery, usePostUserMutation } = usersApi;
