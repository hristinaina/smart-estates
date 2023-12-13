

const CustomDateRangeDialog = ({onConfirm, onCancel}) => {
    return (
        <div id="dialog-overlay">
            <div id="dialog">
                <p id="dialog-title">Add custom date range</p>
                <p id="dialog-message">Choose dates</p>
                <button onClick={onCancel}>CANCEL</button>
                <button onClick={onConfirm}>CONFIRM</button>
            </div>
        </div>
    )
};

export default CustomDateRangeDialog;