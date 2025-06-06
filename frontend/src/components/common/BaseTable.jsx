import React, { useState, useEffect } from 'react';
import {
  DataTable,
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
  Button,
  Tile,
  Tag
} from '@carbon/react';

import { Add } from '@carbon/icons-react';

const transformDataForTable = (data) => {
  return data.map((item, index) => ({
    ...item,
    id: item.id || String(index)
  }))
};

const BaseTable = ({
  title,
  headers,
  rows,
  loading,
  actions = [],
  onAdd,
  renderCustomCell,
  addButtonText = 'Add'
}) => {
  const [sortedRows, setSortedRows] = useState([]);

  useEffect(() => {
    if (rows.length > 0) {
      const defaultSortKey = headers[0]?.key || 'id';
      const sorted = [...rows].sort((a, b) => {
        if (a[defaultSortKey] < b[defaultSortKey]) return -1;
        if (a[defaultSortKey] > b[defaultSortKey]) return 1;
        return 0;
      });
      setSortedRows(sorted);
    }
  }, [rows, headers]);

  return (
    <Tile>
      <div className="cds--data-table-header">
        <h4>{title}</h4>
        {onAdd && <Button onClick={onAdd} renderIcon={Add}>{addButtonText}</Button>}
      </div>
      <DataTable rows={transformDataForTable(sortedRows)} headers={headers} isSortable>
        {({ rows, headers, getHeaderProps, getTableProps, getRowProps }) => {
          const { key: tableKey, ...tableProps } = getTableProps();
          return (
            <Table key={tableKey} {...tableProps}>
              <TableHead>
                <TableRow>
                  {headers.map(header => {
                    const { key, ...headerProps } = getHeaderProps({ header });
                    return (
                      <TableHeader key={header.key} {...headerProps}>
                        {header.header}
                      </TableHeader>
                    );
                  })}
                </TableRow>
              </TableHead>
              <TableBody>
                {rows.map(row => {
                  const { key: rowKey, ...rowProps } = getRowProps({ row });
                  return (
                    <TableRow key={row.id || rowKey} {...rowProps}>
                      {row.cells.map((cell, index) => (
                        <TableCell key={cell.id || `${row.id}-${index}`}>
                          {cell.info.header === 'actions' ? (
                            <div className="flex gap-2">
                              {actions.map((action, actionIndex) => (
                                <Button
                                  key={actionIndex}
                                  kind="ghost"
                                  size="sm"
                                  renderIcon={action.icon}
                                  iconDescription={action.description}
                                  hasIconOnly
                                  onClick={() => action.onClick(row)}
                                />
                              ))}
                            </div>
                          ) : renderCustomCell ? 
                            renderCustomCell(cell, row, index) : 
                            cell.value}
                        </TableCell>
                      ))}
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          );
        }}
      </DataTable>
    </Tile>
  );
};

export default BaseTable;
