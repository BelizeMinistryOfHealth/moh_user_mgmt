import React from 'react';
import { config } from './config';
import * as firebase from 'firebase/app';
import { initializeApp } from 'firebase/app';
import { Provider } from 'react-redux';
import { store } from './store';
import { router } from './routes/router';
import { MantineProvider } from '@mantine/core';
import { RouterProvider } from 'react-router-dom';

let app: firebase.FirebaseApp | null = null;

if (app === null) {
  app = initializeApp(config);
}

function App() {
  return (
    <Provider store={store}>
      <MantineProvider theme={{ colorScheme: 'dark' }}>
        <RouterProvider router={router} />
      </MantineProvider>
    </Provider>
  );
}

export default App;
