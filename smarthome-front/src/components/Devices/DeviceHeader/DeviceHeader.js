import React, { useState } from 'react';
import "./DeviceHeader.css";
import AvailabilityTimeRangeSelector from '../../Dialog/AvailabilityTimeRangeSelector';

const DeviceHeader = ({ handleBackArrow, name }) => {
    const [showDialog, setShowDialog] = useState(false);


    const openTimeRangeDialog = () => {
        console.log("tu sam");
        setShowDialog(true);
    }

    const closeDialog = () => {
        setShowDialog(false);
    }


  return (
    <div>
        <div id="tools">
            <img
                src='/images/arrow.png'
                alt='arrow'
                id='arrow'
                style={{cursor: "pointer" }}
                onClick={handleBackArrow}/>
            <span className='estate-title'>{name}</span>
            <p id="availability" onClick={openTimeRangeDialog}>View device availability</p>
      </div>
      {showDialog && (
                <AvailabilityTimeRangeSelector
                    // onConfirm={this.confirmNewDateRange}
                    onCancel={closeDialog}
                />
        )};
    </div>
  );
};

export default DeviceHeader;
