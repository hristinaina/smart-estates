import React from "react";
import './Dialog.css';

const Dialog = ({ title, message, onConfirm, onCancel}) => {
    return (
        <div id="dialog-overlay">
            <div id="dialog">
                <p id="dialog-title">{title}</p>
                <p id="dialog-message">{message}</p>
                <button onClick={onConfirm}>CONFIRM</button>
                <button onClick={onCancel}>CANCEL</button>
            </div>
        </div>
    )
}

export default Dialog;