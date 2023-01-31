import React from 'react';
import { Button, createStyles } from '@mantine/core';
import * as firebase from 'firebase/auth';
import { useAppDispatch, useTypedSelector } from '../store';
import { logout, selectUser } from '../features/auth/authSlice';
import { Link, useMatch, useNavigate } from 'react-router-dom';

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
  link: {
    textDecoration: 'none',
    color: theme.colors.gray[0],
  },
}));
const Sidebar = () => {
  const { classes } = useStyles();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const authUser = useTypedSelector(selectUser);
  // const matchRoute = useMatchRoute();
  const isLogin = useMatch('/login');
  // const isLogin = matchRoute({ to: '/login' });

  const handleLogout = () => {
    const auth = firebase.getAuth();
    auth.signOut().then(() => {
      dispatch(logout());
      navigate('/login');
    });
  };

  if (!authUser && isLogin) {
    return <></>;
  }
  if (!authUser) {
    return (
      <div className={classes.nav}>
        <Button>
          <Link to={'/login'} className={classes.link} replace>
            Login
          </Link>
        </Button>
      </div>
    );
  }

  return (
    <div className={classes.nav}>
      <Link to={'/users'} className={classes.link}>
        <Button>Users</Button>
      </Link>
      <Button onClick={() => navigate('/users/new')}>Create User</Button>
      <Button onClick={handleLogout}>Logout</Button>
    </div>
  );
};

export default Sidebar;
