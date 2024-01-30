import { Button, Checkbox, Chip, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, FormControlLabel, Snackbar, Table, TableBody, TableCell, TableHead, TablePagination, TableRow, TextField } from "@mui/material";
import { Component } from "react";
import { Navigation } from "../Navigation/Navigation";
import RealEstateService from "../../services/RealEstateService";
import DeviceService from "../../services/DeviceService";
import authService from "../../services/AuthService";
import PermissionService from "../../services/PermissionService";

export class GrantPermission extends Component {
    constructor(props) {
        super(props);

        this.state = {
            realEstate: {},
            devices: [],
            tableData: [],

            emails: [],
            newEmail: '',
            maxEmails: 5,
            selectedDevices: [],
            selectAll: false,
            grantPermissionDisabled: true,

            selectedRows: [],
            page: 0,
            rowsPerPage: 5,
            dialogOpen: false,

            snackbarMessage: '',
            showSnackbar: false,
            open: false,
        };
    }

    async componentDidMount() {
        const valid = authService.validateUser();
        if (!valid) window.location.assign("/");

        const parts = window.location.href.split('/');
        const id = parts[parts.length - 1];

        const realEstate = await RealEstateService.getById(id)
        this.setState({ realEstate: realEstate })

        const devices = await DeviceService.getDevices(id)
        this.setState({ devices: devices})

        const grantedPermissions = await PermissionService.getPermissionsByRealEstateId(id)
        this.setState({ tableData: grantedPermissions})
    }

    // LEFT SIDE
    handleEmailChange = (event) => {
        this.setState({ newEmail: event.target.value });
    }

    handleKeyDown = (event) => {
        const { newEmail, emails, maxEmails } = this.state;

        if (event.key === 'Enter' && newEmail.trim() !== '') {
            if (emails.length < maxEmails) {
                this.setState({ emails: [...emails, newEmail], newEmail: '' }, () => {
                    this.setState({ grantPermissionDisabled: this.checkGrantPermissionDisabled() }); 
                });
            }
        } 
    }

    handleDeleteEmail = (index) => {
        const { emails } = this.state;
        const updatedEmails = [...emails];
        updatedEmails.splice(index, 1);
        this.setState({ emails: updatedEmails });
    }

    handleSelectAll = (event) => {
        const { devices } = this.state;
        const isChecked = event.target.checked;

        if (isChecked) {
            this.setState({ selectedDevices: devices, selectAll: true }, () => {
                this.setState({ grantPermissionDisabled: this.checkGrantPermissionDisabled() }); 
            });
        } else {
            this.setState({ selectedDevices: [], selectAll: false }, () => {
                this.setState({ grantPermissionDisabled: this.checkGrantPermissionDisabled() }); 
            });
        }
    }

    handleDeviceChange = (event, deviceId) => {
        const { selectedDevices } = this.state;
        const isChecked = event.target.checked;

        if (isChecked) {
            this.setState({ selectedDevices: [...selectedDevices, deviceId] }, () => {
                this.setState({ grantPermissionDisabled: this.checkGrantPermissionDisabled() }); 
            });
        } else {
            this.setState({ selectedDevices: selectedDevices.filter(id => id !== deviceId) }, () => {
                this.setState({ grantPermissionDisabled: this.checkGrantPermissionDisabled() }); 
            });
        }
    }

    checkGrantPermissionDisabled = () => {
        const { emails, selectedDevices } = this.state;
        return !(emails.length > 0 && selectedDevices.length > 0);
    };

    handleGrantPermission = async () => {
        const user = authService.getCurrentUser()

        const deviceIds = this.state.selectedDevices.map(device => device.Id);

        const exist = this.isExistSelectedPermissions(deviceIds)

        if (!exist) {
            await PermissionService.sendGrantValues({
                "Emails": this.state.emails,
                "Devices": deviceIds,
                "RealEstateId": this.state.realEstate.Id,
                "RealEstateName": this.state.realEstate.Name,
                "User": user.Name + " " + user.Surname
            })

            this.setState({ selectedDevices: [], selectAll: false, newEmail: "", emails: [], grantPermissionDisabled: true })

            this.setState({ snackbarMessage: "Successfully granted permissions!" });
            this.handleClick();
        }
        else {
            this.setState({ snackbarMessage: "You have already grant permissions. Look at the table!" });
            this.handleClick();
        }
    }

    isExistSelectedPermissions = (deviceIds) => {
        const found = this.state.tableData.some(item =>
            this.state.emails.some(email =>
                deviceIds.some(device =>
                    item.DeviceId === device && item.UserEmail === email
                )
            )
        );
        if (found) {
            console.log("Uslov ispunjen!");
            return true;
        }
        return false;
    }

    // RIGHT SIDE
    handleRowSelect = (event, index) => {
        const { selectedRows } = this.state;
        const selectedIndex = selectedRows.indexOf(index);
        let newSelectedRows = [];
    
        if (selectedIndex === -1) {
            newSelectedRows = newSelectedRows.concat(selectedRows, index);
        } else if (selectedIndex === 0) {
            newSelectedRows = newSelectedRows.concat(selectedRows.slice(1));
        } else if (selectedIndex === selectedRows.length - 1) {
            newSelectedRows = newSelectedRows.concat(selectedRows.slice(0, -1));
        } else if (selectedIndex > 0) {
            newSelectedRows = newSelectedRows.concat(
                selectedRows.slice(0, selectedIndex),
                selectedRows.slice(selectedIndex + 1)
            );
        }
        this.setState({ selectedRows: newSelectedRows });
    };
    
    isSelected = (index) => {
        return this.state.selectedRows.indexOf(index) !== -1;
    };

    handleChangePage = (event, newPage) => {
        console.log(newPage)
        this.setState({ page: newPage })
    };
    
    handleChangeRowsPerPage = (event) => {
        console.log(event.target.value)
        this.setState({ rowsPerPage: event.target.value})
        this.setState({ page: 0 })
    };

    handleDenyPermission = async () => {
        const selectedData = this.state.selectedRows.map(index => this.state.tableData[index]);
        console.log(selectedData)
        await PermissionService.deletePermit(this.state.realEstate.Id, selectedData);

        this.setState({ dialogOpen: false });

        this.setState({ snackbarMessage: "Successfully denied permissions!" });
        this.handleClick();

        const grantedPermissions = await PermissionService.getPermissionsByRealEstateId(this.state.realEstate.Id)
        this.setState({ tableData: grantedPermissions})

        this.setState({ selectedRows: [] })
    };
    
    handleOpenDialog = () => {
        this.setState({ dialogOpen: true });
    };
    
    handleCloseDialog = () => {
        this.setState({ dialogOpen: false });
    };

    // snackbar
    handleClick = () => {
        this.setState({ open: true });
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        this.setState({ open: false });
    };


    render() {
        const { emails, newEmail, selectAll, devices, selectedDevices, selectedRows, dialogOpen, tableData, page, rowsPerPage } = this.state;

        const startIndex = page * rowsPerPage;
        const endIndex = Math.min(startIndex + rowsPerPage, tableData.length);
        const slicedData = tableData.slice(startIndex, endIndex);

        return (
            <div>
                <Navigation />
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Grant new permission</p>
                        <div>
                            <TextField
                                label="Email"
                                value={newEmail}
                                onChange={this.handleEmailChange}
                                onKeyDown={this.handleKeyDown}
                                fullWidth
                                variant="outlined"
                                placeholder="Enter email addresses"
                            />
                            <div>
                                {emails.map((email, index) => (
                                    <Chip
                                        key={index}
                                        label={email}
                                        onDelete={() => this.handleDeleteEmail(index)}
                                        variant="outlined"
                                    />
                                ))}
                            </div>
                        </div>

                        <div style={{ display: 'flex', flexDirection: 'column', marginLeft:"30%", marginTop: "25px", marginBottom: "25px" }}>
                            <FormControlLabel
                                control={
                                    <Checkbox
                                        checked={selectAll}
                                        onChange={(event) => this.handleSelectAll(event)}
                                    />
                                }
                                label={<span style={{ fontWeight: 'bold' }}>Select All</span>}
                            />
                            {devices.map(device => (
                                <FormControlLabel
                                    key={device.Id}
                                    control={
                                        <Checkbox
                                            checked={selectedDevices.includes(device)}
                                            onChange={(event) => this.handleDeviceChange(event, device)}
                                        />
                                    }
                                    label={device.Name}
                                />
                            ))}
                        </div>
                        <Button onClick={this.handleGrantPermission} variant="contained" color="primary" disabled={this.state.grantPermissionDisabled}>Grant Permission</Button>
                    </div>

                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Deny permissions for {this.state.realEstate.Name} </p>
                        {tableData.length == 0 &&(
                            <p>You have not yet granted a permit for the selected real estate</p>
                        )}
                        {tableData.length != 0 && (
                            <div>
                            <Table>
                            <TableHead>
                                <TableRow>
                                <TableCell>User</TableCell>
                                <TableCell>User Email</TableCell>
                                <TableCell>Device Name</TableCell>
                                <TableCell>Delete</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                            {slicedData.map((item, index) => {
                                const isItemSelected = this.isSelected(index);
                                return (
                                    <TableRow
                                    style={{cursor: "pointer"}}
                                    key={index}
                                    selected={isItemSelected}
                                    onClick={(event) => this.handleRowSelect(event, index)}
                                    hover>
                                
                                    <TableCell>{item.User}</TableCell>
                                    <TableCell>{item.UserEmail}</TableCell>
                                    <TableCell>{item.Device}</TableCell>
                                    <TableCell>
                                        <Checkbox
                                        checked={isItemSelected}
                                        onChange={(event) => this.handleRowSelect(event, index)}
                                        />
                                    </TableCell>
                                    </TableRow>
                                );
                                })}
                            </TableBody>
                        </Table>

                        <TablePagination
                            rowsPerPageOptions={[5, 10, 25, 50]}
                            component="div"
                            count={tableData.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            onPageChange={this.handleChangePage}
                            onRowsPerPageChange={this.handleChangeRowsPerPage}
                            />

                        <Button style={{marginTop: "20px"}} onClick={this.handleOpenDialog} variant="contained" color="primary" disabled={selectedRows.length === 0}>Deny Permission</Button>
                        <Dialog
                            open={dialogOpen}
                            onClose={this.handleCloseDialog}
                            aria-labelledby="alert-dialog-title"
                            aria-describedby="alert-dialog-description"
                            >
                            <DialogTitle id="alert-dialog-title">{"Are you sure you want to deny permission?"}</DialogTitle>
                            <DialogContent>
                                <DialogContentText id="alert-dialog-description">
                                You are about to deny permission for {selectedRows.length} selected rows.
                                </DialogContentText>
                            </DialogContent>
                            <DialogActions>
                                <Button onClick={this.handleCloseDialog} color="primary">Cancel</Button>
                                <Button onClick={this.handleDenyPermission} color="primary" variant="contained">Deny</Button>
                            </DialogActions>
                        </Dialog>
                        </div>
                        )}
                    </div>
                </div>
                <Snackbar
                    open={this.state.open}
                    autoHideDuration={2000}
                    onClose={this.handleClose}
                    message={this.state.snackbarMessage}
                    action={this.action}
                />
            </div>
        )
    }
}