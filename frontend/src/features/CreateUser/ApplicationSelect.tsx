import React from 'react';
import { UserApplication } from '../../models/userApplications';
import { Checkbox, Select } from '@mantine/core';
import { UseFormReturnType } from '@mantine/form/lib/types';

type Props = {
  application: UserApplication;
  form: UseFormReturnType<unknown>;
};
const ApplicationSelect = (props: Props) => {
  const { application, form } = props;
  return (
    <div>
      <Checkbox label={application.name} value={application.id} {...form.getInputProps('applicationId')} />
      <Select
        data={application.permissions.map((permission) => {
          return { value: permission, label: permission };
        })}
        label={'Permissions'}
        searchable
        nothingFound={'No permissions found'}
      />
    </div>
  );
};

export default ApplicationSelect;
