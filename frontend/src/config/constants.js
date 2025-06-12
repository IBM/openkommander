export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';
import { Settings, Edit, TrashCan, Add, View, Reset, Group } from '@carbon/icons-react';

export const TABLE_CONFIGS = {
  brokers: {
    title: 'Brokers',
    headers: [
      { key: 'id', header: 'ID' },
      { key: 'host', header: 'Host' },
      { key: 'port', header: 'Port' },
      { key: 'status', header: 'Status' },
      // { key: 'partitionCount', header: 'Partitions' },
      { key: 'actions', header: 'Actions' }
    ],
    actions: [
      {
        icon: Settings,
        description: 'View Config',
        onClick: (row) => row.onViewConfig?.(row)
      }
    ]
  },
  topics: {
    title: 'Topics',
    headers: [
      { key: 'name', header: 'Name' },
      { key: 'partitions', header: 'Partitions' },
      { key: 'replicationFactor', header: 'Replication' },
      { key: 'status', header: 'Status' },
      { key: 'actions', header: 'Actions' }
    ],
    actions: [
      {
        icon: Edit,
        description: 'Edit',
        onClick: (row) => row.onEdit?.(row)
      },
      {
        icon: TrashCan,
        description: 'Delete',
        onClick: (row) => row.onDelete?.(row.name) // Use name instead of id
      }
    ]
  },
  schemas: {
    title: 'Schemas',
    headers: [
      { key: 'subject', header: 'Subject' },
      { key: 'version', header: 'Version' },
      { key: 'compatibility', header: 'Compatibility' },
      { key: 'actions', header: 'Actions' }
    ],
    actions: [
      {
        icon: View,
        description: 'View Versions',
        onClick: (row) => row.onViewVersions?.(row)
      },
      {
        icon: Edit,
        description: 'Edit',
        onClick: (row) => row.onEdit?.(row)
      },
      {
        icon: TrashCan,
        description: 'Delete',
        onClick: (row) => row.onDelete?.(row.subject) // Use subject instead of id
      }
    ]
  },
  consumerGroups: {
    title: 'Consumer Groups',
    headers: [
      { key: 'groupId', header: 'Group ID' },
      { key: 'members', header: 'Members' },
      { key: 'topics', header: 'Topics' },
      { key: 'status', header: 'Status' },
      { key: 'lag', header: 'Lag' },
      { key: 'actions', header: 'Actions' }
    ],
    actions: [
      {
        icon: Group,
        description: 'View Assignments',
        onClick: (row) => row.onViewAssignments?.(row)
      },
      {
        icon: Reset,
        description: 'Reset Offsets',
        onClick: (row) => row.onResetOffsets?.(row.groupId) // Use groupId instead of id
      }
    ]
  },
  acls: {
    title: 'ACL Rules',
    headers: [
      { key: 'principal', header: 'Principal' },
      { key: 'resourceType', header: 'Resource Type' },
      { key: 'resourceName', header: 'Resource Name' },
      { key: 'operation', header: 'Operation' },
      { key: 'permissionType', header: 'Permission' },
      { key: 'actions', header: 'Actions' }
    ],
    actions: [
      {
        icon: Edit,
        description: 'Edit',
        onClick: (row) => row.onEdit?.(row)
      },
      {
        icon: TrashCan,
        description: 'Delete',
        onClick: (row) => row.onDelete?.(row.id)
      }
    ]
  }
};

export const FORM_CONFIGS = {
  topic: {
    initialData: {
      name: '',
      partitions: 1,
      replication_factor: 1, // Changed to snake_case to match API
    }
  },
  brokers: {
    initialData: {
      id: '',
      host: '',
      port: 9092,
      status: '',
      partitionCount: ''
    }
  },
  acl: {
    initialData: {
      principal: '',
      resourceType: 'Topic', 
      resourceName: '',
      operation: 'Read',
      permissionType: 'Allow'
    }
  },
  schema: {
    initialData: {
      subject: '',
      schema: '',
      compatibility: 'BACKWARD'
    }
  }
};
