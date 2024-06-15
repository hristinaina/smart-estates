import React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { TablePagination } from '@mui/material';

class TableOfActions extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            sortOrder: 'asc',
            sortBy: 'User',
            page: 0,
            rowsPerPage: 5,
        };
    }

    getSortedData() {
        const { logData } = this.props;
        const { sortOrder, sortBy } = this.state;

        const sortedKeys = Object.keys(logData).sort((a, b) => {
            const entryA = logData[a][sortBy];
            const entryB = logData[b][sortBy];

            if (sortBy === "Date") {
                const dateA = new Date(a);
                const dateB = new Date(b);

                return sortOrder === 'asc' ? dateA - dateB : dateB - dateA;
            }

            if (typeof entryA === 'string' && typeof entryB === 'string') {
                return sortOrder === 'asc' ? entryA.localeCompare(entryB) : entryB.localeCompare(entryA);
            } else {
                return sortOrder === 'asc' ? entryA - entryB : entryB - entryA;
            }
        });

        return sortedKeys.map((timestamp) => ({
            ...logData[timestamp],
            timestamp: timestamp,
        }));
    }

    handleSort = (columnName) => {
        this.setState((prevState) => ({
            sortOrder: prevState.sortOrder === 'asc' ? 'desc' : 'asc',
            sortBy: columnName,
        }));
    };

    handleChangePage = (event, newPage) => {
        console.log(newPage)
        this.setState({ page: newPage })
    };

    handleChangeRowsPerPage = (event) => {
        console.log(event.target.value)
        this.setState({ rowsPerPage: event.target.value })
        this.setState({ page: 0 })
    };

    formatTimestamp = (timestamp) => {
        const parsedTime = new Date(timestamp);
        return parsedTime.toLocaleString(); // Adjust the format as needed
    };

    render() {
        const { sortOrder, sortBy, rowsPerPage, page } = this.state;
        const sortedData = this.getSortedData();
        console.log(sortedData);
        const startIndex = page * rowsPerPage;
        const endIndex = Math.min(startIndex + rowsPerPage, sortedData.length);
        const slicedData = sortedData.slice(startIndex, endIndex);

        return (
            <div>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('User')}>
                                User {sortBy === 'User' && sortOrder === 'asc' && '↑'}
                                {sortBy === 'User' && sortOrder === 'desc' && '↓'}
                            </TableCell>
                            <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('Action')}>
                                Action {sortBy === 'Action' && sortOrder === 'asc' && '↑'}
                                {sortBy === 'Action' && sortOrder === 'desc' && '↓'}
                            </TableCell>
                            <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('Percentage')}>
                                Percentage {sortBy === 'Percentage' && sortOrder === 'asc' && '↑'}
                                {sortBy === 'Percentage' && sortOrder === 'desc' && '↓'}
                            </TableCell>
                            <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('Plug')}>
                                Plug Id{sortBy === 'Plug' && sortOrder === 'asc' && '↑'}
                                {sortBy === 'Plug' && sortOrder === 'desc' && '↓'}
                            </TableCell>
                            <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('Date')}>
                                Date {sortBy === 'Date' && sortOrder === 'asc' && '↑'}
                                {sortBy === 'Date' && sortOrder === 'desc' && '↓'}
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {slicedData.map((entry) => {
                            return (
                                <TableRow key={entry["timestamp"]}>
                                    <TableCell>{entry.User}</TableCell>
                                    <TableCell>{entry.Action}</TableCell>
                                    <TableCell>{parseInt(entry.Percentage)}%</TableCell>
                                    <TableCell>{entry.Plug != -1 ? entry.Plug + 1 : -1}</TableCell>
                                    <TableCell>{this.formatTimestamp(entry["timestamp"])}</TableCell>
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>

                <TablePagination
                    rowsPerPageOptions={[5, 10, 25, 50]}
                    component="div"
                    count={sortedData.length}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    onPageChange={this.handleChangePage}
                    onRowsPerPageChange={this.handleChangeRowsPerPage}
                />
            </div>
        );
    }
}

export default TableOfActions;
