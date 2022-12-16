import React from 'react';
import { createBrowserRouter, redirect } from 'react-router-dom';
import Root from './Root';
import LoginRoute from './LoginRoute';
import { importFromLocalStorage } from '../localStorage';
import { AuthSlice } from '../features/auth/authSlice';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    loader: () => {
      const localState: AuthSlice = importFromLocalStorage() as AuthSlice;
      console.log({ localState });
      if (!localState || !localState?.user) {
        throw redirect('/login');
      }
      return localState;
    },
  },
  {
    path: '/login',
    element: <LoginRoute />,
  },
]);
