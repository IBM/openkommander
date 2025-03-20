import React from 'react';
import {
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
  Tile,
} from '@carbon/react';

const ClusterMetrics = ({ metrics, loading }) => {
  if (loading) {
    return <div>Loading...</div>;
  }

  const rows = [
    { key: 'brokers', label: 'Brokers', value: metrics.brokers },
    { key: 'topics', label: 'Topics', value: metrics.topics },
    { key: 'health', label: 'Health', value: metrics.health },
    { key: 'messagesPerSecond', label: 'Messages/sec', value: metrics.messagesPerSecond }
  ];

  return (
    <Tile>
      <h4>Cluster Metrics</h4>
      <Table>
        <TableBody>
          {rows.map(row => (
            <TableRow key={row.key}>
              <TableCell>{row.label}</TableCell>
              <TableCell>{row.value}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Tile>
  );
};

export default ClusterMetrics;
