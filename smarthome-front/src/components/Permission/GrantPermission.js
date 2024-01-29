import { Button, Checkbox, Chip, FormControlLabel, Snackbar, TextField } from "@mui/material";
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

            emails: [],
            newEmail: '',
            maxEmails: 5,

            selectedDevices: [],
            selectAll: false,

            grantPermissionDisabled: true,

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
    }

    // RIGHT SIDE
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
        console.log(this.state.emails)
        console.log(deviceIds)
        console.log(this.state.realEstate.Id)
        console.log(this.state.realEstate.Name)
        console.log(user.Name + " " + user.Surname)

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
        const { emails, newEmail, selectAll, devices, selectedDevices } = this.state;

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
                        {/* <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                Email:
                                <select style={{width: "200px", cursor: "pointer"}}
                                    className="new-real-estate-select"
                                    value={pickedValue}
                                    onChange={(e) => this.setState({ pickedValue: e.target.value })}>
                                    <option value={email}>{ email }</option>
                                    <option value="auto">auto</option>
                                    <option value="none">none</option>
                                </select>
                            </label>
                            <label>
                                Start Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={startDate} onChange={(e) => this.setState({ startDate: e.target.value })} />
                            </label>
                            <label>
                                End Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={endDate} onChange={(e) => this.setState({ endDate: e.target.value })} />
                            </label>
                            <br />
                            <Button type="submit" id='sp-data-button'>Filter</Button>
                        </form>
                        <LogTable logData={logData} /> */}
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