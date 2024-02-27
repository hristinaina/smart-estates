import React, {useState} from 'react';
import "./Dialog.css";
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'; 
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { Line } from "react-chartjs-2";


const AvailabilityTimeRangeSelector = ({onConfirm, onCancel}) => {
    const [selectedTimeRange, setSelectedTimeRange] = useState({
        startTime: '',
        endTime: '',
        lastXitem: [],
    });

    const [isGraphVisible, setIsGraphVisible] = useState(false);

    const testData = {
        labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
        datasets: [
          {
            label: 'Monthly Sales',
            data: [65, 59, 80, 81, 56, 55, 40],
            fill: false,
            borderColor: 'rgba(75,192,192,1)',
            borderWidth: 2,
          },
        ],
      };

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

    const goBack = () => {
        setIsGraphVisible(false);
    };

    const confirm = () => {
        setIsGraphVisible(true);
    };

    const close = () => {
        setIsGraphVisible(false);
        onCancel();
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
                                <p><b>OR</b> choose predefined time range:</p>
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
                            <button onClick={close}>CANCEL</button>
                            <button onClick={confirm}>CONFIRM</button>
                         </div>
                    )};
                   
                    {/* <button onClick={() => onConfirm(selectedDate)}>CONFIRM</button> */}
                    {isGraphVisible && (
                        <div id="dialog">
                            <p id="dialog-title">Availability graph</p>
                            <p id='dialog-message'>Start Date: 2023-12-31</p>
                            <p id='dialog-message'>End Date: 2024-01-07</p>
                            <Line data={testData}/>
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