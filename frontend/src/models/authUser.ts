export type AuthUser = {
  uid: string;
  email: string;
  token: string;
  refreshToken: string;
  expires?: number;
};

export const RoleValues = ['AdminRole', 'SrRole', 'PeerNavigatorRole', 'AdherenceCounselorRole'] as const;
export type Role = typeof RoleValues[number];

export const OrgValues = ['NAC', 'MOHW', 'BFLA', 'GoJoven'] as const;
export type Org = typeof OrgValues[number];

export type User = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  role: Role;
  org: Org;
  enabled: boolean;
};
