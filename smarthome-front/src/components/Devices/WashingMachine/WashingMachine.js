import { Component } from "react";
import authService from "../../../services/AuthService";
import DeviceService from "../../../services/DeviceService";
import { FormControl, Snackbar, Switch, TextField, Typography } from "@mui/material";
import LogTable from "../AirConditioner/LogTable";
import { Button, Input } from "reactstrap";
import { Navigation } from "../../Navigation/Navigation";

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
            wmName: ""
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());    
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }
    
    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await DeviceService.getDeviceById(this.id, 'http://localhost:8081/api/wm/');
        console.log(device.Device.Device.Name)
        this.setState({ wmName: device.Device.Device.Name });
        console.log(this.Name)
        // const updatedMode = device.Mode.split(',').map((m) => ({
        //     name: m,
        //     switchOn: false,
        //     temp: 20.0
        // }));
        // this.setState({mode: updatedMode})
        // this.setState({device: device})

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
    
    // componentDidUpdate(prevProps, prevState) {
    // // Logika koja se izvršava nakon ažuriranja props ili stanja
    // }
    
    // componentWillUnmount() {
    // // Logika koja se izvršava pre nego što se komponenta ukloni
    // }
    
    // handleEvent = () => {
    // // Metoda za rukovanje događajem
    // }

    handleBackArrow() {
        window.location.assign("/devices")
    }
    
    render() {
        const { wmName, device, logData, mode, email, startDate, endDate, currentTemp, pickedValue } = this.state;

        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' alt='arrow' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow} />
                <span className='estate-title'>{wmName}</span>
                <div className='sp-container'>
                    <div id="ac-left-card">
                        {/* <p className='sp-card-title'>Supported Modes</p>
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
                        )
                    })} */}

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