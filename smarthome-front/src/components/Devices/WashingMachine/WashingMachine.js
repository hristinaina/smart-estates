import { Component } from "react";
import authService from "../../../services/AuthService";
import DeviceService from "../../../services/DeviceService";
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, FormControl, InputLabel, MenuItem, Select, Snackbar, Switch, TextField, Typography, Button, Table, TableHead,
    TableBody,
    TableRow,
    TableCell,
    IconButton, } from "@mui/material";
import LogTable from "../AirConditioner/LogTable";
import { Navigation } from "../../Navigation/Navigation";
import WashingMachineService from "../../../services/WashingMachineService";
import { Close } from "@mui/icons-material";
import mqtt from 'mqtt';
import PieChart from "../AirConditioner/PieChart";


export class WashingMachine extends Component {
    constructor(props) {
        super(props);
        this.state = {
            device: {},
            mode: [],
            switchOn: false,
            logData: [],
            startDate: '',
            pickedValue: '',
            snackbarMessage: '',
            showSnackbar: false,
            wmName: "",
            isDialogOpen: false,
            isAllScheduledModeOpen: false,
            scheduledModes: [],
            selectedMode: '',
            selectedDateTime: '',
            isSaveDisabled: true,
            user: null,
            open: false,
            remainingTime: "00:00:00",
            intervalId: null,
            email: ""
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

        const user = authService.getCurrentUser();

        const scheduledModes = await WashingMachineService.getScheduledMode(device.Device.Device.Id);

        this.setState({
            wmName: device.Device.Device.Name,
            mode: updatedMode,
            device: device,
            user: user,
            email: user.Name + " " + user.Surname,
            scheduledModes: scheduledModes
        });

        const logData = await WashingMachineService.getWMHistoryData(this.id, 'none', "", "");  
        this.setState({ 
            logData: logData.result,
            pickedValue: "none",
            startDate: "",
            endDate: "",
        });

        try {
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-wm",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('wm/schedule');
                });

                // Handle incoming MQTT messages
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessageForWM(topic, message);
                });
            }
        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessageForWM(topic, message) {
        const result = JSON.parse(message.toString())
        
        if(result.id === this.state.device.Device.Device.Id) {
            const selectedMode = this.state.mode.find(item => item.Id === result.mode);
            this.handleSwitchToggle(selectedMode, "mqtt")
        }    
    }

    sendDataToSimulation = (mode, previous, isSwitchOn, user) => {
        const topic = "wm/switch/" + this.id;

        var message = {
            "Mode": mode,
            "Switch": isSwitchOn,
            "Previous": previous,
            "UserEmail": user,
        }
        this.mqttClient.publish(topic, JSON.stringify(message));
    }

    sendDataToGetScheduleMode = () => {
        const topic = "wm/get/" + this.id;

        var message = {
            "Get": true,
        }
        this.mqttClient.publish(topic, JSON.stringify(message));
    }

    handleSwitchToggle = (selectedItem, source) => {
        const { mode, user } = this.state;
        let userName = user.Name + " " + user.Surname

        if (source === "mqtt") {
            userName = "auto";
        }

        const currentlyActiveMode = mode.find(item => item.switchOn);
    
        if (currentlyActiveMode) {
            if (selectedItem.Name === currentlyActiveMode.Name) {
                // user turn off
                this.sendDataToSimulation(selectedItem.Name, '', !selectedItem.switchOn, userName);
                this.updateRemainingTime(0);
            } else {
                this.sendDataToSimulation(selectedItem.Name, currentlyActiveMode.Name, !selectedItem.switchOn, userName);
                this.updateRemainingTime(selectedItem.Duration);
            }
        } else {
            this.sendDataToSimulation(selectedItem.Name, '', !selectedItem.switchOn, userName);
            this.updateRemainingTime(selectedItem.Duration);
        }
    
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
        this.validateInputs();
        this.setState({ isDialogOpen: true });
    };

    handleCloseDialog = () => {
        this.setState({ isDialogOpen: false });
    };

    handleChange = (event) => {
        const value = event.target.value;
        this.setState({ selectedMode: value }, () => {
            this.validateInputs(); 
        });
    };

    handleChangeDateTime = (event) => {
        const value = event.target.value;
        this.setState({ selectedDateTime: value }, () => {
            this.validateInputs(); 
        });
    };

    validateInputs = () => {
        const { selectedMode, selectedDateTime } = this.state;
        const currentDate = new Date();
        const selectedDateObject = new Date(selectedDateTime);
        const isSaveDisabled = !(selectedMode && selectedDateTime && selectedDateObject > currentDate);
        this.setState({ isSaveDisabled });
    };

    isAlreadyScheduleForSelectedDate = (duration) => {
        const selectedDateTime = new Date(this.state.selectedDateTime);
    
        const programDurationInMilliseconds = duration * 60 * 1000;
    
        const startTimeOfNewProgram = selectedDateTime.getTime();
        const endTimeOfNewProgram = startTimeOfNewProgram + programDurationInMilliseconds;
    
        return this.state.scheduledModes.some(mode => {
            const modeStartTime = new Date(mode.StartTime).getTime();
            let modeDuration = 60
            if(mode.ModeId === 1) modeDuration = 120
            else if (mode.ModeId === 3) modeDuration = 30
            else if (mode.ModeId === 4) modeDuration = 90
            const modeEndTime = modeStartTime + (modeDuration * 60 * 1000); 
    
            if (startTimeOfNewProgram === modeStartTime) {
                return true;
            }
    
            return ((startTimeOfNewProgram < modeEndTime && endTimeOfNewProgram > modeStartTime) || (startTimeOfNewProgram < modeEndTime && startTimeOfNewProgram > modeStartTime) || (startTimeOfNewProgram < modeStartTime && endTimeOfNewProgram > modeStartTime));
        });
    }
    

    handleSave = async () => {
        const selectedMode = this.state.mode.find(mode => mode.Name === this.state.selectedMode);

        const requestData = {
            DeviceId: this.state.device.Device.Device.Id,
            StartTime: this.state.selectedDateTime,
            ModeId: selectedMode.Id           
        };

        if(!this.isAlreadyScheduleForSelectedDate(selectedMode.Duration)) {
            await WashingMachineService.scheduledMode(requestData);

            const scheduledModes = await WashingMachineService.getScheduledMode(this.state.device.Device.Device.Id);
            console.log(scheduledModes)
            this.setState({ scheduledModes: scheduledModes })

            this.sendDataToGetScheduleMode()
        }         
        else {
            this.setState({ snackbarMessage: "Scheduled mode already exists at the selected time" });
            this.handleClick();
        }

        this.handleCloseDialog();
    };

    handleShowScheduleModes = async () => {
        this.setState({ isAllScheduledModeOpen: true });

        const modeIdToNameMap = {};
        this.state.device.Mode.forEach(mode => {
            modeIdToNameMap[mode.Id] = mode.Name;
        });

        const scheduledModesWithNames = this.state.scheduledModes.map(mode => ({
            ...mode,
            ModeName: modeIdToNameMap[mode.ModeId]
        }));

        this.setState({ scheduledModes: scheduledModesWithNames });
    };    
    
    handleFormSubmit = async (e) => {
        e.preventDefault();

        const { startDate, endDate, pickedValue } = this.state;
        if(new Date(startDate) > new Date(endDate)) {
            this.setState({ snackbarMessage: "Start date must be before end date" });
            this.handleClick();
            return 
        }
        const logData = await WashingMachineService.getWMHistoryData(this.id, pickedValue, startDate, endDate);
        this.setState({
            logData: logData.result,
        });
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

    updateRemainingTime = (duration) => {
        const { intervalId } = this.state;

        if (intervalId) {
            clearInterval(intervalId); 
        }
        
        if (duration === 0) {
            console.log("nema nista")
            this.setState({ remainingTime: "00:00:00" }); 
            return;
        }
        
        console.log("ipak ima")
        const durationInSeconds = duration * 60; 
        let remainingTime = durationInSeconds;
        
        const newIntervalId = setInterval(() => {
            remainingTime -= 1;
            if (remainingTime <= 0) {
                clearInterval(newIntervalId); 
                remainingTime = 0;
            }
            this.setState({ remainingTime: this.formatTime(remainingTime), intervalId: newIntervalId }); // Postavi novi intervalId u stanje komponente
        }, 1000);

        this.setState({ intervalId: newIntervalId }); 
    };

    formatTime(seconds) {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const remainingSeconds = seconds % 60;
    
        // Formatiranje vremena
        const formattedHours = String(hours).padStart(2, '0');
        const formattedMinutes = String(minutes).padStart(2, '0');
        const formattedSeconds = String(remainingSeconds).padStart(2, '0');
    
        return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
    }

    
    render() {
        const { wmName, remainingTime, isDialogOpen, selectedMode, scheduledModes, selectedDateTime, logData, mode, email, startDate, endDate, pickedValue } = this.state;

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
                        <div>
                            <p>Remaining Time: {remainingTime}</p>
                        </div>                                               
                        {mode.map((item, index) => {
                        return (
                        <div key={`${item.name}-${index}`} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                            <Typography style={{ fontSize: '1.1em' }}>Off</Typography>
                            <Switch
                                checked={item.switchOn}
                                onChange={() => this.handleSwitchToggle(item, "user")}
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
                            Select mode, enter date and time, then click Save.
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
                                <Button variant="contained" onClick={this.handleSave} disabled={this.state.isSaveDisabled} color="primary">
                                    Save
                                </Button>
                            </div>
                        </DialogActions>
                    </Dialog>

                    <Button id='sp-data-button' style={{marginLeft: "15px"}} onClick={this.handleShowScheduleModes}>Show schedule modes</Button>
                    <Dialog open={this.state.isAllScheduledModeOpen} onClose={() => this.setState({ isAllScheduledModeOpen: false })}>
                        <div style={{ backgroundColor: "white", padding: "20px", width: "400px" }}>
                            <DialogTitle variant="h6">Scheduled Modes</DialogTitle>
                            <IconButton
                                aria-label="close"
                                onClick={() => this.setState({ isAllScheduledModeOpen: false })}
                                sx={{ position: 'absolute', right: 8, top: 8, color: 'gray' }}>
                                <Close />
                            </IconButton>
                            <DialogContent>
                                {scheduledModes.length === 0 ? (
                                    <Typography variant="body1">There are no scheduled modes</Typography>
                                ) : (
                                    <Table>
                                        <TableHead>
                                            <TableRow>
                                                <TableCell>Start Time</TableCell>
                                                <TableCell>Mode</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {scheduledModes.map((mode, index) => (
                                                <TableRow key={index}>
                                                    <TableCell>{mode.StartTime}</TableCell>
                                                    <TableCell>{mode.ModeName}</TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                )}
                            </DialogContent>
                        </div>
                    </Dialog>

                    </div>

                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Switch History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                User:
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

                <div id='statistics'>
                    <p className='sp-card-title'>Statistic</p>
                    <p>Graphs are based on switch history data</p>
                    <div>
                        <p className='sp-card-title'>Mode usage percentage %</p>
                        <PieChart data={logData} graph={1} />

                        <p className='sp-card-title'>Device activity percentage %</p>
                        <PieChart data={logData} graph={2} />

                        <p className='sp-card-title'>User usage percentage %</p>
                        <PieChart data={logData} graph={3} />
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