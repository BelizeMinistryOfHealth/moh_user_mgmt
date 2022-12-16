import React from 'react';
import reactLogo from './assets/react.svg';
import './App.css';
import { config } from './config';
import * as firebase from 'firebase/app';
import { initializeApp } from 'firebase/app';
import { Provider } from 'react-redux';
import { store, useTypedSelector } from './store';
import { RouterProvider } from 'react-router-dom';
import { router } from './routes/router';
import { selectUser } from './features/auth/authSlice';

let app: firebase.FirebaseApp | null = null;

if (!app) {
  app = initializeApp(config);
  //   console.log({ app });
}

function App() {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  );
}

export default App;
