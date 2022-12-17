import React from 'react';
import { createBrowserRouter, redirect } from 'react-router-dom';
import Root from './Root';
import LoginRoute from './LoginRoute';
import { importFromLocalStorage } from '../localStorage';
import { AuthSlice } from '../features/auth/authSlice';
import { STORAGE_KEYS } from '../appConstants';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    loader: (): AuthSlice => {
      const localState: AuthSlice = importFromLocalStorage() as AuthSlice;
      if (!localState || !localState?.user) {
        throw redirect('/login');
      }
      if (localState.user.expires && localState.user.expires > new Date().getTime()) {
        localStorage.removeItem(STORAGE_KEYS.USER_DATA);
        throw redirect('/login');
      }
      return localState;
    },
  },
  {
    path: '/login',
    element: <LoginRoute />,
    loader: (): AuthSlice => {
      const localState: AuthSlice = importFromLocalStorage() as AuthSlice;
      if (localState && localState?.user) {
        throw redirect('/');
      }
      return localState;
    },
  },
]);
