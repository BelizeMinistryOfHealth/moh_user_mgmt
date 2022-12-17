import { createSelector, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState } from '../../store';
import { saveToLocalStorage } from '../../localStorage';
import { AuthUser } from '../../models/authUser';
import { STORAGE_KEYS } from '../../appConstants';
import jwt_decode from 'jwt-decode';

export type AuthSlice = {
  user?: AuthUser | null;
};

const initialState: AuthSlice = {};

const slice = createSlice({
  name: 'authSlice',
  initialState,
  reducers: {
    createUser(state, action: PayloadAction<{ user: AuthUser }>) {
      const { user } = action.payload;
      state.user = user;
      const decoded = jwt_decode(user.token);
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      state.user.expires = decoded.exp;
      saveToLocalStorage(state);
      return state;
    },
    logout(state) {
      localStorage.removeItem(STORAGE_KEYS.USER_DATA);
      state.user = null;
      return state;
    },
  },
});

export const selectUser = createSelector(
  (state: RootState) => state.authSlice,
  (authSlice): AuthUser | undefined | null => authSlice.user,
);

export const { createUser, logout } = slice.actions;

export default slice.reducer;
