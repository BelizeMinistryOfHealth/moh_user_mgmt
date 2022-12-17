export type AuthUser = {
  uid: string;
  email: string;
  token: string;
  refreshToken: string;
  expires?: number;
};

export type User = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
};
