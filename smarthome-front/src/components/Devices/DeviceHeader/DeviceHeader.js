import React from 'react';
import "./DeviceHeader.css";

const DeviceHeader = ({ handleBackArrow, name }) => {
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
            <p id="availability">View device availability</p>
      </div>
    </div>
  );
};

export default DeviceHeader;
