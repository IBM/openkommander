import React, { useEffect } from 'react';
import { Stack } from '@carbon/react';
import BaseTable from './BaseTable';
import BaseModal from './BaseModal';
import ErrorNotification from './ErrorNotification';
import { useResourceManager } from '../../hooks/useResourceManager';

const DefaultFormContent = () => null;

const ResourcePage = ({
  endpoint,
  tableConfig,
  formConfig,
  renderFormContent,
  renderCustomCell,
  transformRow,
  customActions = []
}) => {
  const {
    data,
    loading,
    error,
    setError,
    fetchData,
    isModalOpen,
    setIsModalOpen,
    isEditing,
    formData,
    setFormData,
    handleSubmit,
    handleAdd,
    handleEdit,
    handleDelete
  } = useResourceManager(endpoint, formConfig);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const enhancedRows = data.map(item => ({
    ...item,
    ...transformRow?.(item),
    onEdit: () => handleEdit(item), // Ensure handleEdit is called with the correct item
    onDelete: () => handleDelete(item.id), // Ensure handleDelete is called with the correct ID
    ...customActions.reduce((acc, action) => ({
      ...acc,
      [action.key]: action.handler
    }), {})
  }));

  return (
    <>
      <ErrorNotification error={error} onClose={() => setError(null)} />
      <Stack gap={7}>
        <BaseTable
          {...tableConfig}
          rows={enhancedRows}
          loading={loading}
          renderCustomCell={renderCustomCell}
          onAdd={handleAdd}
        />
      </Stack>
      {/* Only render modal if formConfig is provided */}
      {formConfig && (
        <BaseModal
          isOpen={isModalOpen}
          isEditing={isEditing}
          title={tableConfig.title.slice(0, -1)} // Remove 's' from plural title
          onClose={() => setIsModalOpen(false)}
          onSubmit={handleSubmit}
        >
          {(renderFormContent || DefaultFormContent)({ formData, setFormData, isEditing })}
        </BaseModal>
      )}
    </>
  );
};

export default ResourcePage;
