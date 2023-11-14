
import { Component } from 'react';
import './Devices.css';
import { Navigation } from '../Navigation/Navigation';
import DeviceService from '../../services/DeviceService'

export class Devices extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
          };
    }

    async componentDidMount() {
        try {
          const result = await DeviceService.getDevices(1);
          this.setState({ data: result });
        } catch (error) {
          // Handle error
        }
      }

    render() {
        const { data } = this.state;

        return (
            <div>
                <Navigation />
                <p id="add-real-estate">
                    <img alt="." src="/images/plus.png" id="plus" />
                    Add New Device
                </p>
                <DevicesList devices={data} />
            </div>
        )
    }
}

const DevicesList = ({ devices }) => {
    const chunkSize = 5; // Number of items per row
  
    const chunkArray = (arr, size) => {
      return Array.from({ length: Math.ceil(arr.length / size) }, (v, i) =>
        arr.slice(i * size, i * size + size)
      );
    };
  
    const rows = chunkArray(devices, chunkSize);
  
    return (
      <div id='devices-container'>
        {rows.map((row, rowIndex) => (
          <div key={rowIndex} className='device-row'>
            {row.map((device, index) => (
              <div key={index} className='real-estate-card'>
                <img
                  alt='device'
                  src={device.image} 
                  className='device-img'
                />
                <div className='device-info'>
                  <p className='device-title'>{device.name}</p>
                  <p className='device-text'>{device.type}</p>
                  <p className='device-text state-color'>{device.status}</p>
                </div>
              </div>
            ))}
          </div>
        ))}
      </div>
    );
  };
