import { Component } from "react";
import authService from "../../../services/AuthService";
import { Navigation } from "../../Navigation/Navigation";
import VehicleGateService from "../../../services/VehicleGateService";
import './VehicleGate.css';

export class VehicleGate extends Component {
    constructor(props) {
        super(props);
        this.state = {
            device: {},
            email: '',
            startDate: '',
            endDate: '',
        };
        this.mqttClient = null;
        this.id = parseInt(this.extractDeviceIdFromUrl());
        this.Name = '';
    }

    async componentDidMount() {
        const valid = await authService.validateUser();
        if (!valid) window.location.assign("/");

        const device = await VehicleGateService.get(this.id);

        const user = authService.getCurrentUser();
        this.Name = device.ConsumptionDevice.Device.Name;
        await this.setState({device: device});
    }

    extractDeviceIdFromUrl() {
        const parts = window.location.href.split('/');
        return parts[parts.length - 1];
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleModeChange = async(mode) => {
        let device = this.state.device;
        if (device.Mode != mode) {
            device.Mode = mode;
            if (mode === 0) {
                await VehicleGateService.toPrivate(this.id);
            } else {
                await VehicleGateService.toPublic(this.id);
            }
        }
        await this.setState({device: device})
        
    }

    render() {
        return (
            <div>
                <Navigation />
                <img src='/images/arrow.png' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow}/>
                <span className='estate-title'>{this.Name}</span>
                <div className="sp-container">
                    <div id="sp-left-card">
                        <p className="sp-card-title">Vehicle Gate</p>
                        <p className="sp-data-text">Mode</p>
                        <p className="vg-description">{this.state.device.Mode === 0 ? 'Private' : 'Public'}</p>
                        <img src='/images/private.png' className={`vg-icon ${this.state.device.Mode === 1 ? 'unlocked': ''}`} onClick={ () => this.handleModeChange(0)}/>
                        <img src='/images/public.png' className={`vg-icon ${this.state.device.Mode === 0 ? 'unlocked': ''}`} onClick={ () => this.handleModeChange(1)}/>
                        <p className="sp-data-text">State</p>
                        <p className="vg-description">{this.state.device.IsOpen === true ? 'Opened' : 'Closed'}</p>
                        <img src='/images/closed-gate.png' className={`vg-icon ${this.state.device.IsOpen === true ? 'unlocked' : ''}`} />
                        <img src='/images/opened-gate.png' className={`vg-icon ${this.state.device.IsOpen === false ? 'unlocked' : ''}`} />
                    </div>
                </div>
            </div>
        )
    }
}