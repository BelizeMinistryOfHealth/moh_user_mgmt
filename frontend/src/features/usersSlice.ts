import { User } from '../models/authUser';
import { createSelector, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState } from '../store';

export type UsersSlice = {
  selectedUser: User | null;
};

const initialState: UsersSlice = {
  selectedUser: null,
};

const slice = createSlice({
  name: 'usersSlice',
  initialState,
  reducers: {
    chooseUser(state, action: PayloadAction<{ user: User }>) {
      state.selectedUser = action.payload.user;
    },
  },
});

export const selectChosenUser = createSelector(
  (state: RootState) => state.usersSlice,
  (usersSlice: UsersSlice) => usersSlice.selectedUser,
);

export const { chooseUser } = slice.actions;

export default slice.reducer;
