import React from 'react';
import { useForm } from '@mantine/form';
import { Button, createStyles, Group, PasswordInput, TextInput, Title } from '@mantine/core';
import * as firebase from 'firebase/auth';
import { createUser } from '../features/auth/authSlice';
import { useAppDispatch } from '../store';
import { AuthUser } from '../models/authUser';
import { useNavigate } from 'react-router-dom';

const useStyles = createStyles(() => ({
  login: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'flex-start',
    alignItems: 'center',
    height: '100%',
    gap: '1rem',
  },
}));

const Login = () => {
  const navigate = useNavigate();
  const { classes } = useStyles();
  const form = useForm({
    initialValues: {
      email: '',
      password: '',
      error: '',
    },

    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
    },
  });

  const dispatch = useAppDispatch();

  const handleSubmit = async (input: { email: string; password: string; error: string }) => {
    const auth = firebase.getAuth();
    try {
      const result = await firebase.signInWithEmailAndPassword(auth, input.email, input.password);
      const user: AuthUser = {
        uid: result.user.uid,
        email: result.user.email ?? '',
        token: await result.user.getIdToken(),
        refreshToken: result.user.refreshToken,
      };
      dispatch(createUser({ user }));
      navigate('/');
    } catch (error) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      form.setErrors(error.message);
    }
  };

  return (
    <div className={classes.login}>
      <Title color={'white'} size={'h1'}>
        Log In
      </Title>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <TextInput required withAsterisk label={'Email'} {...form.getInputProps('email')} />
        <PasswordInput withAsterisk required label={'Password'} {...form.getInputProps('password')} />
        <Group position={'right'} mt={'md'}>
          <Button type={'submit'}>Submit</Button>
        </Group>
      </form>
    </div>
  );
};

export default Login;
