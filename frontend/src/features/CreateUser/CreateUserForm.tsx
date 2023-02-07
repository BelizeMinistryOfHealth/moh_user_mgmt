import React from 'react';
import { Button, Group, TextInput, Title } from '@mantine/core';
import { usePostUserMutation } from '../../api/usersApi';
import { useForm } from '@mantine/form';
import { Navigate, useNavigate } from 'react-router-dom';

const CreateUserForm = () => {
  const navigate = useNavigate();
  const [postUser, { isSuccess: saved, isLoading: isSaving }] = usePostUserMutation();
  const form = useForm({
    initialValues: {
      firstName: '',
      lastName: '',
      email: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      firstName: (value) => (value.length < 2 ? 'First name must have at least 2 characters' : null),
      lastName: (value) => (value.length < 2 ? 'First name must have at least 2 characters' : null),
    },
  });

  const handleSubmit = async (input: { firstName: string; lastName: string; email: string }) => {
    form.validate();
    console.log({ input });
    if (form.isValid()) {
      // await postUser({ ...input, userApplications }).unwrap();
    }
  };
  if (isSaving) return <>Saving...</>;
  if (saved) {
    return <Navigate to={'/users'} />;
  }

  return (
    <div>
      <Title color={'white'} size={'h2'}>
        Create User
      </Title>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <TextInput required withAsterisk label={'First Name'} {...form.getInputProps('firstName')} />
        <TextInput required withAsterisk label={'Last Name'} {...form.getInputProps('lastName')} />
        <TextInput required withAsterisk label={'Email'} {...form.getInputProps('email')} />
        <div style={{ marginTop: '2rem' }}>
          <Title color={'white'} size={'h3'}>
            Applications
          </Title>
          {/*<Checkbox.Group label={'Permissions'} orientation={'vertical'} onChange={setPermissions}>*/}
          {/*  {applications &&*/}
          {/*    applications.permissions.map((permission) => (*/}
          {/*      <Checkbox key={permission} label={permission} value={permission} />*/}
          {/*    ))}*/}
          {/*</Checkbox.Group>*/}
        </div>
        <Group mt={'sm'}>
          <Button type={'submit'}>Save</Button>
          <Button type={'button'} variant={'outline'} onClick={() => navigate('/users')}>
            Cancel
          </Button>
        </Group>
      </form>
    </div>
  );
};

export default CreateUserForm;
