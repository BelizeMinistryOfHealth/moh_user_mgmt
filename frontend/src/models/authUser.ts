import { UserApplication } from './userApplications';

export type AuthUser = {
  uid: string;
  email: string;
  token: string;
  refreshToken: string;
  expires?: number;
};

export type RawUser = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  userApplications: UserApplication[] | null | undefined;
};
export type User = Omit<RawUser, 'userApplications'> & {
  userApplications: UserApplication | null;
};
