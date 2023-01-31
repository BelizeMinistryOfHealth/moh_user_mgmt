import React from 'react';
import AppShell from '../components/AppShell';
import { RootLoader } from './router';
import { useAppDispatch, useTypedSelector } from '../store';
import { createUser, selectUser } from '../features/auth/authSlice';
import { Outlet, useLoaderData } from 'react-router-dom';

const Home = () => {
  const data = useLoaderData() as RootLoader;
  const userInLocalStorage = data.user;
  const dispatch = useAppDispatch();
  const user = useTypedSelector(selectUser);
  if (userInLocalStorage && !user) {
    dispatch(createUser({ user: userInLocalStorage }));
  }
  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
};

export default Home;
