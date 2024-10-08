import React, { Component } from 'react';
import '../Devices.css';
import'./AirConditioner.css'
import { Navigation } from '../../Navigation/Navigation';
import authService from '../../../services/AuthService'
import mqtt from 'mqtt';
import { TextField, Button, Snackbar, Typography, Switch, Input, FormControl } from '@mui/material';
import DeviceService from '../../../services/DeviceService';
import LogTable from './LogTable';
import PieChart from './PieChart';
import SpecialModeForm from './SpecialModeForm';
import DeviceHeader from '../DeviceHeader/DeviceHeader';
import PermissionService from '../../../services/PermissionService';
import { ConstructionOutlined } from '@mui/icons-material';


export class AirConditioner extends Component {
    connected = false;

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
            splitMode: [],
            convertedSpecialMode: [],
            userEmails: []
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = "";
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/ac/');
        console.log(device)
        const result = await this.convertFormat(device.SpecialMode)
        this.setState({convertedSpecialMode: result})
        const splitMode = device.Mode.split(',')
        this.setState({splitMode: splitMode})

        const updatedMode = splitMode.map((m) => ({
            name: m,
            switchOn: false,
            temp: 20.0
        }));
        this.setState({mode: updatedMode})
        this.setState({device: device})

        const user = authService.getCurrentUser();
        this.Name = device.Device.Device.Name;
        const logData = await DeviceService.getACHistoryData(this.id, user.Name + " " + user.Surname, "", "");      
        const data = this.setAction(logData.result)
        let users = await PermissionService.getAllUsers(this.id, device.Device.Device.RealEstate);
        users.push("auto");
        users.push("all");
        this.setState({ 
            logData: data,
            email: user.Name + " " + user.Surname,
            pickedValue: user.Name + " " + user.Surname,
            startDate: "",
            endDate: "",
            userEmails: users
        });

        try {
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-ac",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('ac/temp');
                    this.mqttClient.subscribe('ac/action');
                });

                // Handle incoming MQTT messages
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessage(topic, message);
                });
            }
        } catch (error) {
            console.log("Error trying to connect to broker");
            console.log(error);
        }
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    setAction = (data) => {
        for (const key in data) {
            if (data[key].Action === 0) {
                data[key].Action = "Turn off";
            } else if (data[key].Action === 1) {
                data[key].Action = "Turn on";
            }

            if (data[key].User.includes("@")) {
                data[key].User = this.state.email;
            }
        }
        return data
    }

    handleSwitchToggle = (item) => {
        console.log(item)
        console.log("uslo je i ovde")
        let i = 0
        const { mode, email } = this.state;
        let canTurnOn = false
        if(!item.switchOn)
            canTurnOn = this.canTurnOn(item.name, item.temp)
    
        if(canTurnOn || item.switchOn)
        {
            const updatedMode = mode.map((m) => { 
                // ako je bio upaljen, posalji da se gasi   
                if(m.switchOn && item.name !== m.name) {  
                    console.log("dugme prvo")
                    this.sendDataToSimulation(item.name, item.temp, m.name, !item.switchOn, email)
                    ++i                   
                }
                if (!m.switchOn && m.name === item.name) {
                    return {
                        ...m,
                        switchOn: true,
                    };
                } 
                else { 
                    // turn off others                   
                    return {
                        ...m,
                        switchOn: false,
                    };
                }
            });       
            // ovo znaci da nista pre toga nije bilo ukljuceno/iskljuceno
            if(i===0) {
                console.log("dugme drugo")
                this.sendDataToSimulation(item.name, item.temp, '', !item.switchOn, email)
            }                   
            
            this.setState({ mode: updatedMode });
            console.log(updatedMode)
        }
    };

    handleScheduledToggle = (item) => {
        let i = 0
        const { mode } = this.state;

        const updatedMode = mode.map((m) => { 
            // ako je bio upaljen, posalji da se gasi   
            if(m.switchOn && item.name !== m.name) {  
                console.log("zakazano prvo")
                this.sendDataToSimulation(item.name, item.temp, m.name, !item.switchOn, "auto")
                ++i                   
            }
            if (!item.switchOn && m.name === item.name) {
                return {
                    ...m,
                    switchOn: true,
                };
            } 
            else { 
                // turn off others                   
                return {
                    ...m,
                    switchOn: false,
                };
            }
        });       
        // ovo znaci da nista pre toga nije bilo ukljuceno/iskljuceno
        console.log(i)
        console.log(item.switchOn)
        if(i===0 && !item.switchOn) {
            console.log("zakazano drugo")
            console.log(item.name)
            console.log(item.temp)
            console.log(!item.switchOn)
            this.sendDataToSimulation(item.name, item.temp, '', !item.switchOn, "auto")
        }                   
        
        this.setState({ mode: updatedMode });
        console.log(updatedMode)
    };

    sendDataToSimulation = (mode, temp, previous, isSwitchOn, user) => {
        const topic = "ac/switch/" + this.id;

        var message = {
            "Mode": mode,
            "Switch": isSwitchOn,
            "Temp": temp,
            "Previous": previous,
            "UserEmail": user,
        }
        this.mqttClient.publish(topic, JSON.stringify(message));
    }

    canTurnOn = (mode, temp) => {
        const { device, currentTemp } = this.state
        // da li je uneta temperatura u rasponu device.min i device.max
        if(device.MinTemperature > temp || temp > device.MaxTemperature) {
            this.setState({ snackbarMessage: "Temperature out of the range" });
            this.handleClick();
            return false
        } 
        // ako je grejanje ukljuceno da li je veca od trenutne
        else if(mode === "Heating" && temp <= currentTemp) {
            this.setState({ snackbarMessage: "Invalid heating temperature" });
            this.handleClick();
            return false
        }
        // ako je hladjenje ukljuceno da li je manja od trenutne
        else if(mode === "Cooling" && temp >= currentTemp) {
            this.setState({ snackbarMessage: "Invalid cooling temperature" });
            this.handleClick();
            return false
        }
        return true
    } 
    
    handleTemperatureChange = (item, event) => {
        const { mode } = this.state;
        // console.log(event.target.value)
        // event.target.value = event.target.value == '' ? this.state.device.MinTemperature : event.target.value

        const currentIndex = mode.findIndex((m) => m.name === item.name);

        if (currentIndex !== -1) {
            const updatedMode = [...mode];
            updatedMode[currentIndex] = {
                ...updatedMode[currentIndex],
                temp: parseFloat(event.target.value),
            };

            this.setState({ mode: updatedMode });
        }
    }

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        const result = JSON.parse(message.toString())
        if (result.id === this.id) {
            this.setState({
                currentTemp: result.temp
            });
            if(result.mode != null) {
                this.handleScheduledToggle({
                    name: result.mode,
                    switchOn: !result.switch,
                    temp: result.temp,
                })
            }
        }
    }

    handleFormSubmit = async (e) => {
        e.preventDefault();

        const { email, startDate, endDate, pickedValue } = this.state;
        console.log(email, startDate, endDate);
        if(new Date(startDate) > new Date(endDate)) {
            this.setState({ snackbarMessage: "Start date must be before end date" });
            this.handleClick();
            return 
        }
        const logData = await DeviceService.getACHistoryData(this.id, pickedValue, startDate, endDate);
        console.log(logData.result)
        const data = this.setAction(logData.result)
        this.setState({
            logData: data,
        });
    };

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleBackArrow() {
        window.location.assign("/devices")
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

    handleAddSpecialMode = (specialMode) => {
        // todo izmeni ovo i napravi api da se izmeni ovo 
        const updatedSchedule = specialMode.map(item => ({
            ...item,
            temperature: parseFloat(item.temperature)
        }));
        this.setState({ specialModes: updatedSchedule });
        console.log('Adding special mode:', updatedSchedule);
    };

    async convertFormat (data) {
        console.log(data)
        return data.map(item => ({
            start: item.StartTime,
            end: item.EndTime,
            selectedMode: item.Mode,
            temperature: item.Temperature,
            selectedDays: item.SelectedDays.split(",")
        }));
    }   

    render() {
        const { device, logData, mode, email, startDate, endDate, currentTemp, pickedValue, splitMode, convertedSpecialMode, userEmails } = this.state;

        return (
            <div>
                <Navigation />
                <DeviceHeader handleBackArrow={this.handleBackArrow} name={this.Name} />
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Supported Modes</p>
                        <div style={{marginBottom: "25px"}}>
                            <div>
                                <span className='ac-current-temp'>Min temp:  </span>
                                <span><b>{device.MinTemperature}</b></span>
                                <span style={{marginLeft: "50px"}}></span>
                                <span className='ac-current-temp'>Max temp:  </span>
                                <span><b>{device.MaxTemperature}</b></span>
                            </div>
                            <span className='ac-current-temp'>Current temp:  </span>
                            <span><b>{ currentTemp }</b></span>                         
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
                            <span style={{ flex: 1 }}>{item.name}</span>                            
                                <FormControl style={{ width: '80px' }}>
                                {item.name !== 'Ventilation' && (
                                    <Input
                                        type="number"
                                        value={item.temp}
                                        onChange={(event) => this.handleTemperatureChange(item, event)}
                                        inputProps={{
                                            min: device.MinTemperature, 
                                            max: device.MaxTemperature,
                                        }}
                                    />
                                    )}
                                </FormControl>                          
                        </div>
                        )})}
                        <div style={{marginTop: "50px"}}>
                            <SpecialModeForm onAdd={this.handleAddSpecialMode} acModes={splitMode} minTemp={device.MinTemperature} maxTemp={device.MaxTemperature} specialModeFromDevice={convertedSpecialMode} fromDevice={true} id={this.id}/>
                        </div>
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
                                    {userEmails.map(email => (
                                        <option key={email} value={email}>
                                        {email}
                                        </option>
                                    ))}
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
