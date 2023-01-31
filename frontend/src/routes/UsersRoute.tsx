import React from 'react';
import { useGetUsersQuery } from '../api/usersApi';
import UsersTable from '../components/UsersTable';
import { User } from '../models/authUser';
import { useTypedSelector } from '../store';
import { selectUser } from '../features/auth/authSlice';
import { Navigate } from 'react-router-dom';

const UsersRoute = () => {
  const { data, isLoading, isFetching, isError } = useGetUsersQuery();
  const user = useTypedSelector(selectUser);
  if (!user) {
    return <Navigate to={'/login'} replace />;
  }
  if (isLoading || isFetching) {
    return <>Loading...</>;
  }

  if (isError) {
    console.log({ data, isError, user });
    return <>Error</>;
  }

  return <UsersTable users={data as User[]} />;
};

export default UsersRoute;
