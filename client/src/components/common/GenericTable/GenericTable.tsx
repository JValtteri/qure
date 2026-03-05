import './GenericTabel.css';

import { useState, useMemo } from 'react';


type SortOrder = 'asc' | 'desc';

export type GenericTableProps<T> = {
    /** Array of data objects */
    data: T[];

    /** Which keys of T should be shown as columns (in order) */
    columns: (keyof T)[];

    /** Optional key that uniquely identifies a row – defaults to `'id'` if present */
    rowKey?: keyof T;

    /** Called when a row is clicked. Receives the full row object. */
    onRowClick?: (row: T) => void;

    /** Enable filtering by first column (firs column must be string) */
    filterable?: boolean

    /** Enable sorting */
    sortable?: boolean;

    /** Overrides default sort. Called when a column header is clicked – receives the column key */
    onCustomSort?: (column: keyof T) => void;
};

function GenericTable<T>({
    data,
    columns,
    rowKey = 'id' as keyof T,
    onRowClick,
    filterable,
    sortable,
    onCustomSort,
}: GenericTableProps<T>) {
    //const initSortKey =
    const [filterText, setFilterText] = useState<string>('');
    const [sortColumn, setSortColumn] = useState<keyof T>(columns[0]);
    const [sortOrder, setSortOrder] = useState<SortOrder>('asc');

    const getRowKey = (row: T, idx: number) => row[rowKey] !== undefined ? String(row[rowKey]) : idx;
    const isSortedCol = (col: keyof T) => sortColumn === col;

    filterable = filterable && data.length > 0 && typeof data[0][columns[0]] === 'string';

    const defaultSort = (a: T, b: T) => {
        const aVal = a[sortColumn];
        const bVal = b[sortColumn];
        if (aVal === bVal) return 0;
        const comparisonValue =
            typeof aVal === 'string'
                ? (aVal as string).localeCompare(bVal as string)
                : ((aVal as number) < (bVal as number) ? -1 : 1);
        return sortOrder === 'asc' ? comparisonValue : -comparisonValue;
    }

    const filteredAndSorted = useMemo(() => {
        const filtered = filterable
            ? data.filter((l) => (l[columns[0]] as string).toLowerCase().includes(filterText.toLowerCase()))
            : data;
        const sorted = sortable
            ? [...filtered].sort((a, b) => defaultSort(a, b))
            : filtered;
        return sorted;
    }, [data, filterText, sortColumn, sortOrder]);


    const handleSort = (column: keyof T) => {
        if (column === sortColumn) {
            setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
        } else {
            setSortColumn(column);
            setSortOrder('asc');
        }
    };

    const onSort = onCustomSort != null ? onCustomSort : handleSort;


    return (
        <>
            {filterable &&
            <input
                type="text"
                placeholder="Search by name…"
                value={filterText}
                onChange={(e) => setFilterText(e.target.value)}
                className="filterInput"
            />}

            <table className="generic-table">
                <thead>
                    <tr>
                        {columns.map((col) => (
                            <th
                                key={String(col)}
                                className="generic-title-column-header"
                                onClick={onSort ? () => onSort(col) : undefined}
                                style={{ cursor: sortable ? 'pointer' : 'default' }}
                            >
                                {String(col).charAt(0).toUpperCase() + String(col).slice(1)}
                                {isSortedCol(col) && sortOrder && (
                                    <span className="sortIndicator">
                                        {sortOrder === 'asc' ? '▲' : '▼'}
                                    </span>
                                )}
                            </th>
                        ))}
                    </tr>
                </thead>

                <tbody>
                    {filteredAndSorted.length > 0 ? (
                        filteredAndSorted.map((row, idx) => (
                            <tr
                                key={getRowKey(row, idx)}
                                className="row"
                                onClick={() => onRowClick?.(row)}
                            >
                                {columns.map((col) => (
                                    <td key={String(col)} className="generic-td">
                                        {String(row[col])}
                                    </td>
                                ))}
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={columns.length} className="noMatches">
                                No matches
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </>
    );
}

export default GenericTable;
