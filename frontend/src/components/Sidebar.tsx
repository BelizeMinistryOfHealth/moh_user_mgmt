import React from 'react';
import { Button, createStyles, Text } from '@mantine/core';
import * as firebase from 'firebase/auth';
import { useAppDispatch } from '../store';
import { logout } from '../features/auth/authSlice';
import { useNavigate } from 'react-router-dom';

const useStyles = createStyles((theme) => ({
  nav: {
    backgroundColor: theme.colors.gray[9],
    border: `solid 1px ${theme.colors.gray[5]}`,
    display: 'flex',
    gap: '2rem',
    padding: '1rem',
    justifyContent: 'stretch',
    flex: 1,
  },
}));
const Sidebar = () => {
  const { classes } = useStyles();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const handleLogout = () => {
    const auth = firebase.getAuth();
    dispatch(logout());
    auth.signOut().then(() => {
      navigate('/login');
    });
  };
  return (
    <div className={classes.nav}>
      <Button>Users</Button>
      <Button>Create User</Button>
      <Button onClick={handleLogout}>Logout</Button>
    </div>
  );
};

export default Sidebar;
