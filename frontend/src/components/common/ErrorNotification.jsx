import React from 'react';
import { InlineNotification } from '@carbon/react';

const ErrorNotification = ({ error, onClose }) => {
  if (!error) return null;
  
  return (
    <InlineNotification
      kind="error"
      title="Error"
      subtitle={error}
      onClose={onClose}
    />
  );
};

export default ErrorNotification;