import React from 'react';
import { useGetApplicationsQuery, useGetUserQuery } from '../api/usersApi';
import { createLoader, withLoader } from '@ryfylke-react/rtk-query-loader';
import { Button, Checkbox, createStyles, Group, TextInput, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useNavigate, useParams } from 'react-router-dom';
import { User } from '../models/authUser';

const useStyles = createStyles(() => ({
  editPage: {
    display: 'flex',
    gap: '1rem',
    flexDirection: 'column',
  },
}));

const editFormLoader = createLoader({
  queries: () => {
    const params = useParams();
    if (!params.userId) throw new Error('No user id provided');
    const user = useGetUserQuery(params.userId);
    const application = useGetApplicationsQuery();
    return [user, application] as const;
  },
  onLoading: () => <>Loading....</>,
});
const UserEditForm = withLoader((_, queries) => {
  const [permissions, setPermissions] = React.useState<string[]>([]);
  const user = queries[0].data;
  const application = queries[1].data;
  const { classes } = useStyles();
  const navigate = useNavigate();
  const form = useForm({
    initialValues: {
      firstName: user.firstName,
      lastName: user.lastName,
      email: user.email,
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
    },
  });

  const handleSubmit = async (input: Omit<User, 'id' | 'userApplications'>) => {
    form.validate();
    console.log({ input, permissions });
  };

  return (
    <div className={classes.editPage}>
      <Title color={'white'} size={'h2'}>
        Edit Form
      </Title>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <TextInput required withAsterisk label={'First Name'} {...form.getInputProps('firstName')} />
        <TextInput required withAsterisk label={'Last Name'} {...form.getInputProps('lastName')} />
        <TextInput required withAsterisk label={'Email'} {...form.getInputProps('email')} />
        <Checkbox.Group label={'Permissions'} orientation={'vertical'} onChange={setPermissions}>
          {application &&
            application.permissions.map((permission) => (
              <Checkbox label={permission} value={permission} key={permission} />
            ))}
        </Checkbox.Group>
        <Group mt={'sm'}>
          <Button type={'submit'}>Save</Button>
          <Button type={'button'} variant={'outline'} onClick={() => navigate('/users')}>
            Cancel
          </Button>
        </Group>
      </form>
    </div>
  );
}, editFormLoader);

export default UserEditForm;
