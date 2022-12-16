import { configureStore } from '@reduxjs/toolkit';
import authSlice from './features/auth/authSlice';
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

export const createStore = () => {
  return configureStore({
    reducer: {
      authSlice,
    },
  });
};

export const store = createStore();

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export type RootState = ReturnType<typeof store.getState>;
export const useTypedSelector: TypedUseSelectorHook<RootState> = useSelector;
