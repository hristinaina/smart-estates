import { Component } from "react";
import authService from "../../../services/AuthService";
import DeviceService from "../../../services/DeviceService";
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, FormControl, InputLabel, MenuItem, Select, Snackbar, Switch, TextField, Typography, Button } from "@mui/material";
import LogTable from "../AirConditioner/LogTable";
import { Navigation } from "../../Navigation/Navigation";
import WashingMachineService from "../../../services/WashingMachineService";

export class WashingMachine extends Component {
    constructor(props) {
        super(props);
        this.state = {
            device: {},
            mode: [],
            switchOn: false,
            logData: [],
            email: '',
            startDate: '',
            endDate: '',
            pickedValue: '',
            snackbarMessage: '',
            showSnackbar: false,
            open: false,
            temp: 20.0,
            currentTemp: "Loading...",
            wmName: "",
            isDialogOpen: false,
            selectedMode: '',
            selectedDateTime: '',
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl()); 
        this.currentDateTime = new Date().toISOString().slice(0, 16); 
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }
    
    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/wm/');
        
        // add switchOn attribute
        const updatedMode = device.Mode.map(modeItem => ({
            ...modeItem,
            switchOn: false 
        }));

        this.setState({
        wmName: device.Device.Device.Name,
        mode: updatedMode,
        device: device
        });

        // const user = authService.getCurrentUser();
        // this.Name = device.Device.Device.Name;
        // const logData = await DeviceService.getACHistoryData(this.id, 'none', "", "");      
        // const data = this.setAction(logData.result)
        // this.setState({ 
        //     logData: data,
        //     email: user.Email,
        //     pickedValue: "none",
        //     startDate: "",
        //     endDate: "",
        // });

        // try {
        //     if (!this.connected) {
        //         this.connected = true;
        //         this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
        //             clientId: "react-front-nvt-2023-ac",
        //             clean: false,
        //             keepalive: 60
        //         });

        //         // Subscribe to the MQTT topic
        //         this.mqttClient.on('connect', () => {
        //             this.mqttClient.subscribe('ac/temp');
        //             this.mqttClient.subscribe('ac/action');
        //         });

        //         // Handle incoming MQTT messages
        //         this.mqttClient.on('message', (topic, message) => {
        //             this.handleMqttMessage(topic, message);
        //         });
        //     }
        // } catch (error) {
        //     console.log("Error trying to connect to broker");
        //     console.log(error);
        // }
    }

    handleSwitchToggle = (selectedItem) => {
        this.setState(prevState => ({
            mode: prevState.mode.map(item => ({
                ...item,
                switchOn: item === selectedItem ? !selectedItem.switchOn : false
            }))
        }));
    };  

    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleOpenDialog = () => {
        this.setState({ isDialogOpen: true });
    };

    handleCloseDialog = () => {
        this.setState({ isDialogOpen: false });
    };

    handleChange = (event) => {
        this.setState({ selectedMode: event.target.value })
    };

    handleChangeDateTime = (event) => {
        const selectedDateTime = event.target.value;
        this.setState({ selectedDateTime });
    };

    handleSave = async () => {
        const selectedMode = this.state.mode.find(mode => mode.Name === this.state.selectedMode);

        const requestData = {
            DeviceId: this.state.device.Device.Device.Id,
            StartTime: this.state.selectedDateTime,
            ModeId: selectedMode.Id           
        };

        await WashingMachineService.scheduledMode(requestData);

        this.handleCloseDialog();
    };
    
    
    render() {
        const { wmName, isDialogOpen, selectedMode, selectedDateTime, device, logData, mode, email, startDate, endDate, currentTemp, pickedValue } = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' alt='arrow' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{wmName}</span>
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Supported Modes</p>
                        <div style={{marginBottom: "25px"}}>                       
                        </div>                                                 
                        {mode.map((item, index) => {
                        return (
                        <div key={`${item.name}-${index}`} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                            <Typography style={{ fontSize: '1.1em' }}>Off</Typography>
                            <Switch
                                checked={item.switchOn}
                                onChange={() => this.handleSwitchToggle(item)}
                            />
                            <Typography style={{ fontSize: '1.1em' }}>On</Typography>
                            <span style={{ flex: 1, fontWeight: "600", marginLeft: "15px" }}>{item.Name}</span> 
                            <span style={{ flex: 1 }}>{item.Duration} min</span> 
                            <span style={{ flex: 1 }}>{item.Temperature}</span>                          
                        </div>
                        )
                    })}

                    <Button id='sp-data-button' onClick={this.handleOpenDialog}>Initiate Mode</Button>

                    <Dialog open={isDialogOpen} onClose={this.handleCloseDialog}>
                        <DialogTitle>Initiate Mode</DialogTitle>
                        <DialogContent>
                        <DialogContentText>
                            Select cycle, enter date and time, then click Save.
                        </DialogContentText>
                        <FormControl fullWidth sx={{ my: 1 }}>
                            <InputLabel id="cycle-select-label">Select Mode</InputLabel>
                            <Select labelId="cycle-select-label" id="cycle-select" value={selectedMode} onChange={this.handleChange}>
                                {mode.map((item) => (
                                    <MenuItem key={item.Name} value={item.Name}>
                                        {item.Name}
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                        <TextField
                            fullWidth
                            margin="normal"
                            label="Enter Date and Time"
                            type="datetime-local"
                            InputLabelProps={{ shrink: true }}
                            inputProps={{ min: this.currentDateTime }}
                            value={selectedDateTime} 
                            onChange={this.handleChangeDateTime} 
                        />
                        </DialogContent>
                        <DialogActions>
                        <div>
                    <Button onClick={this.handleCloseDialog} color="primary">
                        Close
                    </Button>
                </div>
                <div>
                    <Button variant="contained" onClick={this.handleSave} color="primary">
                        Save
                    </Button>
                </div>
                        {/* <Button onClick={this.handleCloseDialog}>Cancel</Button>
                        <Button onClick={this.handleSave}>Save</Button> */}
                        </DialogActions>
                    </Dialog>
                    </div>

                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Switch History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
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
                        <LogTable logData={logData} />
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