import React, { Component } from "react";
import './Sprinkler.css';
import { Navigation } from "../../Navigation/Navigation";
import { IconButton, Switch, Table, TableCell, TableContainer, TableRow, Typography, Paper, TableBody, Chip, TableHead, Button, TextField, Snackbar } from "@mui/material";
import AddSprinklerSpecialMode from "./AddSprinklerSpecialMode";
import CloseIcon from '@mui/icons-material/Close';
import LogTable from "../AirConditioner/LogTable";
import SprinklerService from "../../../services/SprinklerService";
import mqtt from 'mqtt';
import authService from "../../../services/AuthService";
import DeviceHeader from "../DeviceHeader/DeviceHeader";
import PieChart from "../AirConditioner/PieChart";
import PermissionService from "../../../services/PermissionService";


export class Sprinkler extends Component {
    connected = false;

    constructor(props) {
        super(props);

        this.state = {
            specialModes: [],
            allSpecialModes: [],
            startDate: '',
            endDate: '',
            pickedValue: '',
            email: '',
            logData: [],
            switchOn: false,
            open: false,
            snackbarMessage: '',
            showSnackbar: false,
            username: '',
            userEmails: []
        };
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = '';
        this.mqttClient = null;
    }

    async componentDidMount() {
        const user = authService.getCurrentUser();
        const device = await SprinklerService.get(this.id);
        console.log(device);
        let users = await PermissionService.getAllUsers(this.id, device.ConsumptionDevice.Device.RealEstate);
        console.log(users)
        users.push("auto");
        users.push("all");
        this.setState({ 
            userEmails: users,
            email: user.Name + " " + user.Surname,
        });
        this.Name = device.ConsumptionDevice.Device.Name;
        console.log(this.Name);
        await this.setState({email: user.Email, username: user.Name + " " + user.Surname});   
        console.log('*****************');
        const res = await SprinklerService.getSpecialModes(this.id);
        if (res !== null) {
            let specials = [];
            res.forEach(element => {
                const specialMode = {
                    start: element.StartTime,
                    end: element.EndTime,
                    selectedDays: this.getSelectedDays(element.SelectedDays),
                };
                specials.push(specialMode);
            });
           
            await this.setState({specialModes: specials});
        }
        
        const sprinkler = await SprinklerService.get(this.id);
        if (sprinkler.IsOn == true){
            await this.setState({switchOn: true});
        }

        try {        
            // mqtt connection
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                        clientId: "react-front-nvt-2023-sprinkler",
                        clean: false,
                        keepalive: 60
                });
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('sprinkler/on/' + this.id);
                    this.mqttClient.subscribe('sprinkler/off/' + this.id);
                });
    
                this.mqttClient.on('message', (topic, message) => {
                    if (topic === 'sprinkler/on/' + this.id) {
                        this.handleMqttMessage(topic);
                    } else if (topic === 'sprinkler/off/' + this.id) {
                        this.handleMqttOffMessage(topic);
                    }
                });
            }

        } catch (error) {
            console.error(error);
        }
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    async handleMqttMessage(topic, message) {
        var parts = topic.split("/");
        var lastPart = parts[parts.length - 1];
        var parsedNumber = parseInt(lastPart, 10);
        if (parsedNumber != this.id) {
            return;
        }
        await this.setState({switchOn: true});
    }

    async handleMqttOffMessage(topic, message) {
        var parts = topic.split("/");
        var lastPart = parts[parts.length - 1];
        var parsedNumber = parseInt(lastPart, 10);
        if (parsedNumber != this.id) {
            return;
        }
        await this.setState({switchOn: false});
    }

    getSelectedDays(selectedDays) {
        return selectedDays.split(',').filter(day => day !== "");
    }


    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleAddSpecialMode = (newSpecialModes) => {
        this.setState(prevState => ({
            specialModes: [...prevState.specialModes, ...newSpecialModes],
          }));
    }

    handleFormSubmit = async(e) => {
        e.preventDefault();

        let { email, startDate, endDate, pickedValue } = this.state;
        console.log(pickedValue);
        if(new Date(startDate) > new Date(endDate)) {
            this.setState({ snackbarMessage: "Start date must be before end date" });
            this.handleClick();
            return 
        }
        if (pickedValue == ''){
            pickedValue = email;
        }
        const logData = await SprinklerService.getHistoryData(this.id, pickedValue, startDate, endDate);
        for (const timestamp in logData.result) {
            if (logData.result[timestamp].User == this.state.email) {
                logData.result[timestamp].User = this.state.username;
            }
        }
        console.log('------------');
        console.log(logData.result);
        // const data = this.setAction(logData.result)
        this.setState({
            logData: logData.result,
        });
    }

    handleDelete = async(index) => {
        const res = await SprinklerService.getSpecialModes(this.id);
        let specials = [];
        if (res !== null) {
            res.forEach(element => {
                const specialMode = {
                    id: element.Id,
                    start: element.StartTime,
                    end: element.EndTime,
                    selectedDays: this.getSelectedDays(element.SelectedDays),
                };
                specials.push(specialMode);
            });
        }
        await SprinklerService.deleteMode(specials[index].id);
        await this.setState(prevState => ({
            specialModes: prevState.specialModes.filter((_, ind) => ind !== index),
          }));
    }

    handleSwitchToggle = async() => {
        await this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));

        await SprinklerService.changeState(this.id, this.state.switchOn);
    }

    render() {
        return (
            <div>
                <Navigation/>
                <DeviceHeader handleBackArrow={this.handleBackArrow} name={this.Name} />
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Details</p>
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Typography style={{ fontSize: '1.1em'}}>Off</Typography>
                                <Switch
                                    checked={this.state.switchOn}
                                    onChange={() => this.handleSwitchToggle()}
                                />
                            <Typography style={{ fontSize: '1.1em' }}>On</Typography>
                        </div>
                        {/* <p id="special-mode">Add special mode</p> */}
                        <AddSprinklerSpecialMode
                            onAdd={this.handleAddSpecialMode} 
                            modes={this.state.specialModes}/>
                         <TableContainer component={Paper}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Start</TableCell>
                                        <TableCell>End</TableCell>
                                        <TableCell>Day</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {this.state.specialModes.map((item, index) => (
                                        <TableRow key={index}>
                                            <TableCell>{item.start}</TableCell>
                                            <TableCell>{item.end}</TableCell>
                                            <TableCell>
                                                {item.selectedDays.map((day, dayIndex) => (
                                                    <Chip key={dayIndex} label={day} />
                                                ))}
                                            </TableCell>
                                            <TableCell>
                                                <IconButton color="secondary" onClick={() => this.handleDelete(index)}>
                                                    <CloseIcon />
                                                </IconButton>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Switch History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                Email:
                                <select style={{width: "200px", cursor: "pointer"}}
                                    className="new-real-estate-select"
                                    value={this.state.pickedValue}
                                    onChange={(e) => this.setState({ pickedValue: e.target.value })}>
                                    {this.state.userEmails.map(email => (
                                        <option key={email} value={email}>
                                        {email}
                                        </option>
                                    ))}
                                </select>
                            </label>
                            <label>
                                Start Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={this.state.startDate} onChange={(e) => this.setState({ startDate: e.target.value })} />
                            </label>
                            <label>
                                End Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={this.state.endDate} onChange={(e) => this.setState({ endDate: e.target.value })} />
                            </label>
                            <br />
                            <Button type="submit" id='sp-data-button'>Filter</Button>
                            </form>
                            <LogTable logData={this.state.logData} hide={true}/>
                        </div>
                    </div>
                    <div id='statistics'>
                    <p className='sp-card-title'>Statistic</p>
                    <p>Graphs are based on switch history data</p>
                    <div>
                        <p className='sp-card-title'>Device activity percentage %</p>
                        <PieChart data={this.state.logData} graph={4} />

                        <p className='sp-card-title'>User usage percentage %</p>
                        <PieChart data={this.state.logData} graph={3} />
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

    handleClick = () => {
        this.setState({ open: true });
    };

    handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        this.setState({ open: false });
    };

}