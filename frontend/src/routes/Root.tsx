import React from 'react';
import reactLogo from '../assets/react.svg';
import { useAppDispatch, useTypedSelector } from '../store';
import { createUser, selectUser, User } from '../features/auth/authSlice';
import { importFromLocalStorage, saveToLocalStorage } from '../localStorage';

const Root = () => {
  const user = useTypedSelector(selectUser);
  const [fetchState, setFetchState] = React.useState('start');
  const dispatch = useAppDispatch();
  const stateInLocalStorage = importFromLocalStorage();
  if (!user && stateInLocalStorage?.user) {
    console.log({ stateInLocalStorage });
    dispatch(createUser({ user: stateInLocalStorage.user as User }));
  }
  console.log({ user });
  React.useEffect(() => {
    if (user?.token && fetchState === 'start') {
      setFetchState('fetching');
      fetch('https://users-mgmt-e46d3zpgka-ue.a.run.app/users', {
        headers: {
          Authorization: `Bearer ${user?.token}`,
        },
      }).then((result) => {
        console.log({ result });
        setFetchState('stop');
      });
    }
  }, [fetchState]);
  return (
    <div className="App">
      <div>
        <a href="https://vitejs.dev" target="_blank" rel="noreferrer">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
        <a href="https://reactjs.org" target="_blank" rel="noreferrer">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <p className="read-the-docs">Click on the Vite and React logos to learn more</p>
    </div>
  );
};

export default Root;
