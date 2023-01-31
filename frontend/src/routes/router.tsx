import React from 'react';
import LoginRoute from './LoginRoute';
import { importFromLocalStorage } from '../localStorage';
import { createUser, logout } from '../features/auth/authSlice';
import { STORAGE_KEYS } from '../appConstants';
import Home from './Home';
import UsersRoute from './UsersRoute';
import { createBrowserRouter } from 'react-router-dom';
import { store } from '../store';
import { AuthUser } from '../models/authUser';
import UserEditForm from '../components/UserEditForm';
import CreateUserForm from '../features/CreateUser/CreateUserForm';

export type RootLoader = { user: AuthUser | null };
export const router = createBrowserRouter([
  {
    path: '/',
    element: <Home />,
    loader: (): RootLoader => {
      const localState = importFromLocalStorage();
      if (!localState || !localState.user) {
        store.dispatch(logout);
        return { user: null };
      }
      if (localState && localState.user?.expires && localState.user.expires > new Date().getTime()) {
        store.dispatch(createUser({ user: localState.user }));
        localStorage.removeItem(STORAGE_KEYS.USER_DATA);
        return { user: null };
      }
      if (localState && localState.user) {
        store.dispatch(createUser({ user: localState.user }));
        localStorage.removeItem(STORAGE_KEYS.USER_DATA);
        return { user: null };
      }
      return { user: localState.user };
    },
    children: [
      {
        path: 'users',
        element: <UsersRoute />,
        children: [],
      },
      {
        path: 'users/new',
        element: <CreateUserForm />,
      },
      {
        path: 'users/:userId',
        element: <UserEditForm />,
      },
      {
        path: 'login',
        element: <LoginRoute />,
      },
    ],
  },
]);
