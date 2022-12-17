import React from 'react';
import { useAppDispatch } from '../store';
import { AuthSlice, createUser } from '../features/auth/authSlice';
import { createStyles, Text } from '@mantine/core';
import Sidebar from '../components/Sidebar';
import { useGetUsersQuery } from '../api/usersApi';
import { AuthUser, User } from '../models/authUser';
import { useLoaderData } from 'react-router-dom';
import UsersTable from '../components/UsersTable';

const useStyles = createStyles(() => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    flexDirection: 'column',
    height: '100%',
    gap: '1rem',
  },
  header: {
    display: 'flex',
  },
  main: {
    flex: 3,
    height: '100%',
    marginTop: '1rem',
  },
}));
const Root = () => {
  const dispatch = useAppDispatch();
  const { classes } = useStyles();
  const { data, isLoading, isFetching, isError, isSuccess } = useGetUsersQuery();
  const loaderData = useLoaderData() as AuthSlice;
  if (isLoading || isFetching) {
    dispatch(createUser({ user: loaderData.user as AuthUser }));
    return <>Loading...</>;
  }
  if (isError) {
    return <>Error</>;
  }

  return (
    <div className={classes.root}>
      <Text color={'white'} size={'lg'}>
        MOH EPI User Mgmt
      </Text>
      <div className={classes.header}>
        <Sidebar />
      </div>
      <main className={classes.main}>
        <UsersTable users={data as User[]} />
      </main>
    </div>
  );
};

export default Root;
