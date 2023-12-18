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

  // convertDate = (timestampString) => {
  //   // const timestampString = "2023-12-18 13:58:16 CET";

  //   // Razdvajanje datuma i vremena od ostatka stringa
  //   console.log(timestampString)
  //   try {
  //     const [datePart, timePart, timeZonePart] = timestampString.split(' ');
  //   const [year, month, day] = datePart.split('-');
  //   const [hour, minute, second] = timePart.split(':');

  //   // Konstrukcija objekta Date sa dodatkom vremenske zone
  //   const timestamp = new Date(Date.UTC(year, month - 1, day, hour, minute, second));
  //   // const timeZone = timeZonePart.replace('CET', 'UTC'); // Prilagodavanje formata za parsiranje

  //   // Postavljanje vremenske zone
  //   // timestamp.setUTCHours(timestamp.getUTCHours() + parseInt(timeZone, 10));

  //   console.log(timestamp); 
  //     return timestamp
  //   } catch {

  //   }

  // }

  // compareDate = (date1, date2) => {
  //   const { sortOrder } = this.state;
  //   if(sortOrder == 'asc') {
  //     if(date1 > date2)
  //       return date1
  //     else 
  //       return date2
  //   } else {
  //     if(date1 < date2)
  //       return date1
  //     else 
  //       return date2
  //   }
  // }

  handleChangePage = (event, newPage) => {
    console.log(newPage)
    this.setState({ page: newPage })
    // setPage(newPage);
  };

  handleChangeRowsPerPage = (event) => {
    console.log(event.target.value)
    this.setState({ rowsPerPage: event.target.value})
    // setRowsPerPage(parseInt(event.target.value, 10));
    this.setState({ page: 0 })
    // setPage(0);
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
        rowsPerPageOptions={[5, 10, 25, 50]}  // Prilagodite opcije prema vašim potrebama
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
