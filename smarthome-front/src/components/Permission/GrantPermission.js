import { Button, Checkbox, Chip, FormControlLabel, Paper, Snackbar, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TextField } from "@mui/material";
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
        console.log(grantedPermissions)
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
        // console.log(this.state.emails)
        // console.log(deviceIds)
        // console.log(this.state.realEstate.Id)
        // console.log(this.state.realEstate.Name)
        // console.log(user.Name + " " + user.Surname)

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
        const { emails, newEmail, selectAll, devices, selectedDevices, tableData } = this.state;

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
                        <p className='sp-card-title'>All permissions for {this.state.realEstate.Name} </p>
                        <TableContainer >
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
                            {tableData.map((item, index) => {
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
                        </TableContainer>
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