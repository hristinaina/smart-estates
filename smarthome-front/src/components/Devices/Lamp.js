
import { Component } from 'react';
import {Line} from 'react-chartjs-2';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import mqtt from 'mqtt';
import Switch from '@mui/material/Switch';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import authService from '../../services/AuthService';
import 'chart.js/auto';
import LampService from '../../services/LampService';
import CustomDateRangeDialog from '../Dialog/CustomDateRangeDialog';
import { containerClasses } from '@mui/material';
import DeviceService from '../../services/DeviceService';
import DeviceHeader from './DeviceHeader/DeviceHeader';


export class Lamp extends Component {
    connected = false;

    constructor(props) {
        super(props);
        this.state = {
            device: {},
            lamp: {},
            switchOn: false,
            showCustomDateRangeDialog: false,
            data: {
                labels: [],
                datasets: [
                  {
                    label: 'Lightning',
                    data: [],
                    borderColor: 'rgba(128,104,148,1)',
                    borderWidth: 2,
                    fill: false,
                  },
                ],
            },
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());

        this.options = {
            scales: {
              y: {
                beginAtZero: false,
              },
            },
        };
    
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        // const { device } = this.state;  // todo instead of this get device from back by deviceId
        const device = await DeviceService.get(this.id);
        const lamp = await LampService.get(this.id);
        await this.setState({device: device, lamp: lamp});
        
        const updatedData =
        {
            ...device,
            Value: "Loading...",
        }
        this.setState({
            device: updatedData,
        });

        try {
            // populating graph data
            let data = await LampService.getAllGraphData(this.id);
            let datasets = [];
            data.forEach((value, key) => {
                datasets.push({
                            label: key,
                            data: value,
                            borderColor: LampService.getRandomColor(),
                            borderWidth: 2,
                            fill: false,
                          },)
            });
            let keys = [];
            for (let i = 0; i < 101; i++) {
                keys.push(i);
            }
            await this.setState({ data: {
                labels: keys,
                datasets: datasets
            }});
            // ...
            if (!this.connected) {
                this.connected = true;
                this.mqttClient = mqtt.connect('ws://localhost:9001/mqtt', {
                    clientId: "react-front-nvt-2023-lamp",
                    clean: false,
                    keepalive: 60
                });

                // Subscribe to the MQTT topic for device status
                this.mqttClient.on('connect', () => {
                    this.mqttClient.subscribe('device/data/' + this.id);
                });

                // Handle incoming MQTT messages
                this.mqttClient.on('message', (topic, message) => {
                    this.handleMqttMessage(topic, message);
                });

                if (lamp.IsOn) {
                    this.turnOnSwitch();
                }
                this.startSimulation();
            }
        } catch (error) {
            console.error(error);
        }
    }

    componentWillUnmount() {
        // Disconnect MQTT client on component unmount
        if (this.mqttClient) {
            this.mqttClient.end();
        }
    }

    turnOnSwitch = () => {
        this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));
        
    }

    startSimulation = () => {
        console.log("usaooo");
        const topic = "lamp/switch/" + this.id;
        const message = (true).toString();
        this.mqttClient.publish(topic, message);
    }

    handleSwitchToggle = async() => {
        const topic = "lamp/switch/" + this.id;

        this.setState((prevState) => ({
            switchOn: !prevState.switchOn,
        }));
        // const message = (true).toString();
        // this.mqttClient.publish(topic, message);
        // console.log(message);
        // TODO: check this
        if (this.state.lamp.IsOn == false) {
            console.log("ukljuceno");
            let l = await LampService.turnOn(this.id);
            console.log(l);
            await this.setState({lamp: l});
        } else {
            console.log("iskljuceno");
            let l = await LampService.turnOff(this.id);
            console.log(l);
            await this.setState({lamp: l});
        }
    };

    // Handle incoming MQTT messages
    handleMqttMessage(topic, message) {
        const { device } = this.state;
        const newValue = message.toString();
        let lastIndex = topic.lastIndexOf("/");
        let resultSubstring = topic.substring(lastIndex + 1);
        let curr_id = parseInt(resultSubstring, 10);
        const updatedData =
        {
            ...device,
            Value: newValue + "%",
        }
        if (this.id == curr_id) {
            this.setState({
                device: updatedData,
            });
        }
        
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    openDialog = () => {
        this.setState({showCustomDateRangeDialog: true,})
    }

    closeDialog = () => {
        this.setState({showCustomDateRangeDialog: false,})
    }

    confirmNewDateRange = async (from, to) => {
        from = from.toISOString();
        to = to.toISOString();
        console.log("confirmed");
        console.log(from);
        console.log(to);
        this.closeDialog();

        let customData = await LampService.getGraphData(this.id, from, to);
        let newData = { ...this.state.data }; // shallow copy
        newData.datasets.push({
            label: 'from ' + from.slice(0, 10) + ' to ' + to.slice(0, 10),
            data: LampService.createGraphData(customData.data),
            borderColor: LampService.getRandomColor(),
            borderWidth: 2,
            fill: false,
        });
        await this.setState({data: newData});
    }

    render() {
        const { device, switchOn, lamp } = this.state;

        return (
            <div>
                <Navigation />
                <DeviceHeader handleBackArrow={this.handleBackArrow} name={this.state.device.Name}/>
                <div className='sp-container'>
                    <div id='sp-left-card'>
                        <p className='sp-card-title' style={{marginTop:"3.5em"}}>Id: {this.id}</p>
                        {/* {switchOn ? (<p className='device-text'>Value: {device.Value}</p>) : null} */}
                        <p className='sp-data-text'>Read Illumination: {device.Value}</p>
                        <p className='sp-data-text'>Light bulb is {switchOn && Number(device.Value.slice(0, -1)) < 50 ? 'ON' : 'OFF'}</p>
                        { console.log("Switchhh: ") }
                        {console.log(switchOn) }
                        {console.log(device.Value)}
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Typography>Off</Typography>
                            <Switch
                                checked={switchOn}
                                onChange={this.handleSwitchToggle}
                            />
                            <Typography>On</Typography>
                        </div>
                    </div>
                    <div id='sp-right-card'>
                        <Line key={JSON.stringify(this.state.data)}  data={this.state.data} options={this.options} />
                        <div id='custom-date-range-container'>
                            <button id='custom-date-range' onClick={this.openDialog}>Add custom date range</button>
                        </div>
                    </div>
                </div>

                {this.state.showCustomDateRangeDialog && (
                <CustomDateRangeDialog
                    onConfirm={this.confirmNewDateRange}
                    onCancel={this.closeDialog}
                />
                )}
                
            </div>
        )
    }
}
