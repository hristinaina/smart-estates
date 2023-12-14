import React, { useState } from 'react';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';


const CustomDateRangeDialog = ({onConfirm, onCancel}) => {
    const [displayError, setDisplayError] = useState(false);  // true so that is not displayed at the beginning
    const [errorMessage, setErrorMessage] = useState('Please choose dates');
    const [startDate, setStartDate] = useState(null);
    const [endDate, setEndDate] = useState(null);

    const handleStartDateChange = (date) => {
        setStartDate(date);
        if (endDate != null) {
            setDisplayError(false);
            isDateRangeValid();
        }
    };

    const handleEndDateChange = (date) => {
        setEndDate(date);
        if (startDate != null) {
            setDisplayError(false);
            isDateRangeValid();
        }
    };

    const handleConfirm = () => {
        if (startDate != null && endDate != null && isDateRangeValid()) {
            onConfirm(startDate, endDate);
        } else {
            setDisplayError(true);
        }
    }

    const isDateRangeValid = () => {
        const maxDaysDifference = 30;
        const millisecondsInDay = 24 * 60 * 60 * 1000;
    
        if (startDate && endDate) {
            const daysDifference = Math.floor((endDate - startDate) / millisecondsInDay);
            const today = new Date();
            today.setHours(0, 0, 0, 0); // Set time to midnight
    
            if (startDate > today || endDate > today) {
                setDisplayError(true);
                setErrorMessage('Can\'t choose dates in the future.');
                console.log("returned false0");
            } else if (daysDifference > maxDaysDifference) {
                setDisplayError(true);
                setErrorMessage('Please choose a range within 30 days.');
                console.log("returned false");
                return false;
            } else if (daysDifference < 0) {
                setDisplayError(true);
                setErrorMessage('Please check order of dates.');
                console.log("returned false 2");
                return false;
            } else {
                setDisplayError(false);
                setErrorMessage('');
                console.log("returned true");
                return true;
            }
        }
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
                {displayError && <p id='please-choose-dates'>{errorMessage}</p>}
                <button onClick={onCancel}>CANCEL</button>
                <button onClick={handleConfirm}>CONFIRM</button>
                {/* <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button> */}
                </div>
            </div>
        </LocalizationProvider>
    );
};

export default CustomDateRangeDialog;