import React from 'react';
import {
  Modal,
  Form,
  Stack
} from '@carbon/react';

const BaseModal = ({
  isOpen,
  isEditing,
  title,
  children,
  onClose,
  onSubmit,
  size = 'md'
}) => {
  return (
    <Modal
      open={isOpen}
      modalHeading={isEditing ? `Edit ${title}` : `Add ${title}`}
      primaryButtonText={isEditing ? 'Save Changes' : 'Add'}
      secondaryButtonText="Cancel"
      onRequestClose={onClose}
      onRequestSubmit={onSubmit}
      size={size}
    >
      <Form>
        <Stack gap={7}>
          {children}
        </Stack>
      </Form>
    </Modal>
  );
};

export default BaseModal;
