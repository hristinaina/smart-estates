import React, {useState} from 'react';
import "./Dialog.css";
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';


const AvailabilityTimeRangeSelector = ({onConfirm, onCancel}) => {
    const [selectedTimeRange, setSelectedTimeRange] = useState({
        startTime: '',
        endTime: '',
        lastXitem: [],
    });

    // fn to handle changes in the input fields
    const handleTimeRangeChange = (field, value) => {
        setSelectedTimeRange(prevState => ({
        ...prevState,
        [field]: value,
        }));
    };

    // fn to handle changes in the checkboxes for last X item
    const handleLastXitemChange = item => {
        setSelectedTimeRange(prevState => ({
        ...prevState,
        lastXitem: item,
        }));
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns}>
            <div id="dialog-overlay">
                <div id="dialog">
                    <p id="dialog-title">Add custom date range</p>
                    <p id="dialog-message">Choose starting and ending point</p>
                    <DatePicker
                        className='picker'
                        label="Starting point"
                        // value={startDate}
                        // onChange={handleStartDateChange}
                    />
                    <p></p>
                    <DatePicker
                        className='picker'
                        label="Ending point"
                        // value={endDate}
                        // onChange={handleEndDateChange}
                    />
                    <p></p>
                    {/* {displayError && <p id='please-choose-dates'>{errorMessage}</p>} */}
                    
                    <div>
                        <p>Choose predefined time range:</p>
                        {[6, 12, 24, '1w', '1m'].map(item => (
                            <label key={item} style={{marginRight: '10px'}}>
                                <input
                                    type="checkbox"
                                    checked={selectedTimeRange.lastXitem.includes(item)}
                                    onChange={() => handleLastXitemChange(
                                        selectedTimeRange.lastXitem.includes(item)
                                            ? selectedTimeRange.lastXitem.filter(h => h !== item)
                                            : [...selectedTimeRange.lastXitem, item]
                                    )}
                                />
                                {typeof item === 'number' ? `${item}h` : item}
                            </label>
                        ))}
                    </div>
                    <br/>
                    <button onClick={onCancel}>CANCEL</button>
                    <button >CONFIRM</button>
                    {/* <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button> */}
                </div>
            </div>
        </LocalizationProvider>
    );
}

export default AvailabilityTimeRangeSelector;