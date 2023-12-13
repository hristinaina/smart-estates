import React, { useState } from 'react';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';


const CustomDateRangeDialog = ({onConfirm, onCancel}) => {
    const [isChoosen, setIsChoosen] = useState(true);  // true so that is not displayed at the beginning
    const [startDate, setStartDate] = useState(null);
    const [endDate, setEndDate] = useState(null);

    const handleStartDateChange = (date) => {
        setStartDate(date);
        if (endDate != null) {
            setIsChoosen(true);
        }
    };

    const handleEndDateChange = (date) => {
        setEndDate(date);
        if (startDate != null) {
            setIsChoosen(true);
        }
    };

    const handleConfirm = () => {
        if (startDate != null && endDate != null) {
            onConfirm(startDate, endDate);
        } else {
            setIsChoosen(false);
        }
    }

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns}>
            <div id="dialog-overlay">
                <div id="dialog">
                <p id="dialog-title">Add custom date range</p>
                <p id="dialog-message">Choose starting and ending point</p>
                <DatePicker
                    className='picker'
                    label="Starting point"
                    value={startDate}
                    onChange={handleStartDateChange}
                />
                <p></p>
                <DatePicker
                    className='picker'
                    label="Ending point"
                    value={endDate}
                    onChange={handleEndDateChange}
                />
                <p></p>
                {!isChoosen && <p id='please-choose-dates'>Please choose dates</p>}
                <button onClick={onCancel}>CANCEL</button>
                <button onClick={handleConfirm}>CONFIRM</button>
                {/* <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button> */}
                </div>
            </div>
        </LocalizationProvider>
    );
};

export default CustomDateRangeDialog;