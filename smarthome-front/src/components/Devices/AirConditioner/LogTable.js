import React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { TablePagination } from '@mui/material';

class LogTable extends React.Component {
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
      if(sortBy === "Date") {
        const dateA = new Date(a);
        const dateB = new Date(b);

        if (sortOrder === 'asc') {
          return dateA - dateB;
        }
        else {
          return dateB - dateA;
        }
      }

      if (sortOrder === 'asc') {
        return entryA && entryB ? entryA.localeCompare(entryB) : 0;
      } else {
        return entryA && entryB ? entryB.localeCompare(entryA) : 0;
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
    this.setState({ page: newPage })
  };

  handleChangeRowsPerPage = (event) => {
    console.log(event.target.value)
    this.setState({ rowsPerPage: event.target.value})
    this.setState({ page: 0 })
  };

  render() {
    const { sortOrder, sortBy, rowsPerPage, page } = this.state;
    const sortedData = this.getSortedData();

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
              <TableCell style={{ cursor: "pointer" }} onClick={() => this.handleSort('Mode')}>
                Mode {sortBy === 'Mode' && sortOrder === 'asc' && '↑'}
                {sortBy === 'Mode' && sortOrder === 'desc' && '↓'}
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
                  <TableCell>{entry.Mode}</TableCell>
                  <TableCell>{entry["timestamp"]}</TableCell>
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

export default LogTable;
