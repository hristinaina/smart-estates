import React, {useState} from "react";
import './Dialog.css';

const Dialog = ({ title, message, onConfirm, onCancel, isDiscard}) => {
    const [reason, setReason] = useState('');

    const handleConfirm = () => {
        onConfirm(reason);
    }

    return (
        <div id="dialog-overlay">
            <div id="dialog">
                <p id="dialog-title">{title}</p>
                <p id="dialog-message">{message}</p>
                {isDiscard && (
                    <textarea 
                        id="reason" 
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        placeholder="Write a reason here..."></textarea>
                )}
                <button onClick={handleConfirm}>CONFIRM</button>
                <button onClick={onCancel}>CANCEL</button>
            </div>
        </div>
    )
}

export default Dialog;