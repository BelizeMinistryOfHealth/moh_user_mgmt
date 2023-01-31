import React, { ReactNode } from 'react';
import { createStyles, Text } from '@mantine/core';
import Sidebar from './Sidebar';

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
  nav: {
    display: 'flex',
  },
  main: {
    flex: 3,
    height: '100%',
    marginTop: '1rem',
  },
  header: {
    display: 'flex',
    justifyContent: 'center',
  },
}));

type Props = {
  children: ReactNode;
};
const AppShell = (props: Props) => {
  const { children } = props;
  const { classes } = useStyles();
  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <Text color={'white'} size={'lg'}>
          MOH EPI User Mgmt
        </Text>
      </div>
      <div className={classes.nav}>{<Sidebar />}</div>
      <main className={classes.main}>{children}</main>
    </div>
  );
};

export default AppShell;
