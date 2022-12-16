import { AuthSlice } from './features/auth/authSlice';
import { STORAGE_KEYS } from './appConstants';

export const importFromLocalStorage = () => {
  let savedState = null;
  try {
    savedState = localStorage.getItem(STORAGE_KEYS.USER_DATA);
    console.log({ savedState });
  } catch (error) {
    console.error(error);
  }

  if (savedState) {
    return JSON.parse(savedState) as AuthSlice;
  }

  return savedState as AuthSlice;
};

export const saveToLocalStorage = (state: AuthSlice) => {
  try {
    localStorage.setItem(STORAGE_KEYS.USER_DATA, JSON.stringify(state));
  } catch (error) {
    console.error(error);
  }
};
