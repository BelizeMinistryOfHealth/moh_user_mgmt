import { createSelector, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState } from '../../store';
import { saveToLocalStorage } from '../../localStorage';

export type User = {
  uid: string;
  email: string;
  token: string;
  refreshToken: string;
};

export type AuthSlice = {
  user?: User;
};

const initialState: AuthSlice = {};

const slice = createSlice({
  name: 'authSlice',
  initialState,
  reducers: {
    createUser(state, action: PayloadAction<{ user: User }>) {
      const { user } = action.payload;
      state.user = user;
      saveToLocalStorage(state);
      return state;
    },
  },
});

export const selectUser = createSelector(
  (state: RootState) => state.authSlice,
  (authSlice): User | undefined => authSlice.user,
);

export const { createUser } = slice.actions;

export default slice.reducer;
