import React, { useState, useEffect } from 'react';
import {
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
  Tile,
} from '@carbon/react';

const ClusterMetrics = ({ metrics, loading, lastRefreshed }) => {
  const rows = [
    { key: 'brokers', label: 'Brokers', value: metrics.brokers },
    { key: 'topics', label: 'Topics', value: metrics.topics },
    { key: 'health', label: 'Health', value: metrics.health },
    { key: 'messagesPerMinute', label: 'Messages/Min', value: metrics.messagesPerMinute }
  ];

  const loadingPlaceholderStyle = {
    width: '60px',
    height: '16px',
    backgroundColor: '#e0e0e0',
    borderRadius: '4px',
    filter: 'blur(1px)',
    animation: 'pulse 1.5s ease-in-out infinite'
  };

  // Add the keyframe animation to the document head if it doesn't exist
  React.useEffect(() => {
    const styleId = 'cluster-metrics-animations';
    if (!document.getElementById(styleId)) {
      const style = document.createElement('style');
      style.id = styleId;
      style.textContent = `
        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.5; }
        }
      `;
      document.head.appendChild(style);
    }
  }, []);

  const LoadingPlaceholder = () => (
    <div style={loadingPlaceholderStyle} />
  );

  return (
    <Tile>
      <h4>Cluster Metrics</h4>
      <Table>
        <TableBody>
          {rows.map(row => (
            <TableRow key={row.key}>
              <TableCell>{row.label}</TableCell>
              <TableCell>
                {loading ? <LoadingPlaceholder /> : row.value}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <div style={{ marginTop: '1rem', fontStyle: 'italic', color: '#6f6f6f' }}>
        Last Refreshed: {lastRefreshed.toLocaleTimeString()}
      </div>
    </Tile>
  );
};

export default ClusterMetrics;
