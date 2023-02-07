import React from 'react';
import { Button, Group, Radio, TextInput, Title } from '@mantine/core';
import { usePostUserMutation } from '../../api/usersApi';
import { useForm } from '@mantine/form';
import { Navigate, useNavigate } from 'react-router-dom';
import { Org, OrgValues, Role, RoleValues } from '../../models/authUser';

const CreateUserForm = () => {
  const navigate = useNavigate();
  const [org, setOrg] = React.useState<Org | null>(null);
  const [role, setRole] = React.useState<Role | null>(null);
  const [postUser, { isSuccess: saved, isLoading: isSaving }] = usePostUserMutation();
  const form = useForm({
    initialValues: {
      firstName: '',
      lastName: '',
      email: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      firstName: (value) => (value.length < 2 ? 'First Name must have at least 2 characters' : null),
      lastName: (value) => (value.length < 2 ? 'Last Name must have at least 2 characters' : null),
    },
  });

  const handleSubmit = async (input: { firstName: string; lastName: string; email: string }) => {
    form.validate();
    console.log({ input, org, role });
    if (form.isValid() && org && role) {
      await postUser({ ...input, org, role, enabled: true }).unwrap();
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
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)' }}>
          <div style={{ marginTop: '2rem' }}>
            <Title color={'white'} size={'h3'}>
              Organizations
            </Title>
            <Radio.Group
              orientation={'vertical'}
              onChange={(o: Org) => {
                setOrg(o);
              }}
              withAsterisk
            >
              {OrgValues.map((org) => (
                <Radio key={org} label={org} value={org} />
              ))}
            </Radio.Group>
          </div>
          <div style={{ marginTop: '2rem' }}>
            <Title color={'white'} size={'h3'}>
              Roles
            </Title>
            <Radio.Group
              orientation={'vertical'}
              onChange={(r: Role) => {
                setRole(r);
              }}
            >
              {RoleValues.map((role) => (
                <Radio key={role} label={role} value={role} />
              ))}
            </Radio.Group>
          </div>
        </div>
        <Group mt={'sm'}>
          <Button type={'submit'} disabled={!form.isValid() || role === null || org === null}>
            Save
          </Button>
          <Button type={'button'} variant={'outline'} onClick={() => navigate('/users')}>
            Cancel
          </Button>
        </Group>
      </form>
    </div>
  );
};

export default CreateUserForm;
