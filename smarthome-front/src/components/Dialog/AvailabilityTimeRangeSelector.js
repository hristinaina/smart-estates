import React, {useState} from 'react';
import "./Dialog.css";
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { Line } from "react-chartjs-2";
import DeviceAvailabilityService from '../../services/DeviceAvailabilityService';
import { format } from 'date-fns';


const AvailabilityTimeRangeSelector = ({onConfirm, onCancel}) => {
    const [selectedTimeRange, setSelectedTimeRange] = useState({
        startTime: '',
        endTime: '',
        lastXitem: null,
    });

    const [isGraphVisible, setIsGraphVisible] = useState(false);

    const [chartData, setChartData] = useState({
        labels: [],
        datasets: [
            {
                label: 'Device availability',
                data: [],
                fill: false,
                borderColor: 'rgba(75,192,192,1)',
                borderWidth: 2,
            },
        ],
    });

    const fetchNewData = async (labels, data) => {
        // Replace this with your actual data fetching logic
        setChartData(prevState => ({
            ...prevState,
            labels: labels,
            datasets: prevState.datasets.map(dataset => ({
                ...dataset,
                data: data,
            })),
        }));
    };

    const getDeviceId = () => {
        const parts = window.location.href.split('/');
        return parseInt(parts[parts.length - 1], 10);
    }

    // fn to handle changes in the input fields
    const handleTimeRangeChange = (field, value) => {
        setSelectedTimeRange(prevState => ({
        ...prevState,
        [field]: value ? value : null,
        }));
    };

    // fn to handle changes in the checkboxes for last X item
    const handleLastXitemChange = item => {
        setSelectedTimeRange(prevState => ({
        ...prevState,
        lastXitem: item,
        }));
    };

    const goBack = () => {
        handleTimeRangeChange("startTime", '');
        handleTimeRangeChange("endTime", '');
        setIsGraphVisible(false);
    };

    const confirm = async() => {
        if ((selectedTimeRange.startTime === null || selectedTimeRange.startTime === '') || 
        (selectedTimeRange.endTime === null || selectedTimeRange.endTime === '')) {
            let data = await DeviceAvailabilityService.get(getDeviceId(), "-" + selectedTimeRange.lastXitem, "-1");
            let labels = formatLables(Object.keys(data));
            fetchNewData(labels, Object.values(data));
        } else {
            let data = await DeviceAvailabilityService.get(getDeviceId(), selectedTimeRange.startTime.toISOString(), 
                                                           selectedTimeRange.endTime.toISOString());
            let labels = formatLables(Object.keys(data));
            fetchNewData(labels, Object.values(data));
        }


        setIsGraphVisible(true);
    };

    const close = () => {
        setIsGraphVisible(false);
        onCancel();
    }

    const formatLables = (labels) => {
        let formattedLabels = [];
        Object.keys(labels).forEach(element => {
            console.log("ELEMENT " + labels[element]);
            const date = new Date(labels[element]);
            formattedLabels.push(format(date, "dd.MM.yyyy. hh:mm a"));
        });
        return formattedLabels;
    } 

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns}>
            <div id="dialog-overlay">
                
                    {!isGraphVisible && (
                        <div id="dialog">
                            <p id="dialog-title">Add custom date range</p>
                            <p id="dialog-message">Choose starting and ending point</p>
                            <DatePicker
                                className='picker'
                                label="Starting point"
                                value={selectedTimeRange.startDate || null}
                                onChange={(date) => handleTimeRangeChange('startTime', date)}
                                format="MM/dd/yyyy"
                            />
                            <p></p>
                            <DatePicker
                                className='picker'
                                label="Ending point"
                                value={selectedTimeRange.endDate || null}
                                onChange={(date) => handleTimeRangeChange('endTime', date)}
                                format="MM/dd/yyyy"
                            />
                            <p></p>
                            {/* {displayError && <p id='please-choose-dates'>{errorMessage}</p>} */}
                            
                            <div>
                                <p><b>OR</b> choose predefined time range:</p>
                                {['6h', '12h', '24h', '1w', '30d'].map(item => (
                                    <label key={item} style={{marginRight: '10px'}}>
                                        <input
                                            type="radio"
                                            name='lastXitem'
                                            checked={selectedTimeRange.lastXitem === item}
                                            onChange={() => handleLastXitemChange(item)
                                            }
                                        />
                                        {typeof item === 'number' ? `${item}h` : item}
                                    </label>
                                ))}
                            </div>
                            <br/>
                            <button onClick={close}>CANCEL</button>
                            <button onClick={confirm}>CONFIRM</button>
                         </div>
                    )};
                   
                    {/* <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button> */}
                    {isGraphVisible && (
                        <div id="dialog" >
                            <p id="dialog-title">Availability graph</p>
                            <div> 
                                {((selectedTimeRange.startTime === null || selectedTimeRange.startTime === '') || 
                                (selectedTimeRange.endTime === null || selectedTimeRange.endTime === '')) ? (
                                    <p>
                                        {selectedTimeRange.lastXitem === null ? (
                                            <span>No time data chosen. Go back and choose time.</span>
                                        ) : (
                                            <span>In the last: {selectedTimeRange.lastXitem}</span>
                                        )}
                                    </p>
                                ) : (
                                <React.Fragment>
                                    <p id='dialog-message'>Start Date: {format(selectedTimeRange.startTime, 'MMMM do, yyyy')}</p>
                                    <p id='dialog-message'>End Date: {format(selectedTimeRange.endTime, 'MMMM do, yyyy')}</p>
                                </React.Fragment>
                                )}
                            </div>
                            <div style={{ width: '800px', height: '400px' }}>
                            <Line data={chartData}/>
                            </div>
                            <br/>
                            <button onClick={close}>CLOSE</button>
                            <button onClick={goBack}>GO BACK</button>
                        </div>
                    )};
                </div>
        </LocalizationProvider>
    );
}

export default AvailabilityTimeRangeSelector;