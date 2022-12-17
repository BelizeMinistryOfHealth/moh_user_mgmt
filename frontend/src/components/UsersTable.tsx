import React from 'react';
import { User } from '../models/authUser';
import { Table } from '@mantine/core';

type Props = {
  users: User[];
};
const UsersTable = (props: Props) => {
  const { users } = props;
  const rows = users.map((user) => (
    <tr key={user.id}>
      <td>{user.id}</td>
      <td>{user.firstName}</td>
      <td>{user.lastName}</td>
      <td>{user.email}</td>
    </tr>
  ));
  return (
    <Table striped withBorder>
      <thead>
        <tr>
          <th>ID</th>
          <th>First Name</th>
          <th>Last Name</th>
          <th>Email</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  );
};

export default UsersTable;
