import React from 'react';
import { createStyles } from '@mantine/core';
import { Outlet } from 'react-router-dom';

const useStyles = createStyles(() => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    flexDirection: 'column',
    height: '100%',
    gap: '1rem',
    padding: '2rem 9rem',
  },
  main: {
    flex: 3,
    height: '100%',
    marginTop: '1rem',
  },
}));
const Root = () => {
  const { classes } = useStyles();
  console.log('Root');

  return (
    <div className={classes.root}>
      <main className={classes.main}>
        <Outlet />
      </main>
    </div>
  );
};

export default Root;
