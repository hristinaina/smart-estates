import React, { useState } from 'react';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';


const CustomDateRangeDialog = ({onConfirm, onCancel}) => {
    const [selectedDate, handleDateChange] = useState(null);

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns}>
            <div id="dialog-overlay">
                <div id="dialog">
                <p id="dialog-title">Add custom date range</p>
                <p id="dialog-message">Choose starting and ending point</p>
                <DatePicker
                    className='picker'
                    label="Starting point"
                    value={selectedDate}
                    onChange={(firstDate) => handleDateChange(firstDate)}
                />
                <p></p>
                <DatePicker
                    className='picker'
                    label="Ending point"
                    value={selectedDate}
                    onChange={(secondDate) => handleDateChange(secondDate)}
                />
                <p></p>
                <button onClick={onCancel}>CANCEL</button>
                <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button>
                </div>
            </div>
        </LocalizationProvider>
    );
};

export default CustomDateRangeDialog;