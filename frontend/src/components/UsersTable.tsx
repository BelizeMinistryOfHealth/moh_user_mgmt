import React from 'react';
import { User } from '../models/authUser';
import { createStyles, Table, Title } from '@mantine/core';
import { useAppDispatch } from '../store';
import { chooseUser } from '../features/usersSlice';
import { useNavigate } from 'react-router-dom';

const useStyles = createStyles((theme) => ({
  usersContainer: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1.5rem',
    backgroundColor: theme.colors.gray,
    border: `solid 1px ${theme.colors.gray[5]}`,
    padding: '1rem',
  },
}));

type Props = {
  users: User[];
};
const UsersTable = (props: Props) => {
  const { users } = props;
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { classes } = useStyles();
  return (
    <div className={classes.usersContainer}>
      <Title color={'white'} size={'h2'}>
        Users
      </Title>
      <Table striped withBorder highlightOnHover>
        <thead style={{ backgroundColor: 'gray' }}>
          <tr>
            <th>ID</th>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Email</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr
              key={user.id}
              onClick={() => {
                dispatch(chooseUser({ user }));
                navigate(`/users/${user.id}`);
              }}
            >
              <td>{user.id}</td>
              <td>{user.firstName}</td>
              <td>{user.lastName}</td>
              <td>{user.email}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  );
};

export default UsersTable;
