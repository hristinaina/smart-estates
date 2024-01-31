import React, { Component } from "react";
import './Sprinkler.css';
import { Navigation } from "../../Navigation/Navigation";
import { IconButton, Switch, Table, TableCell, TableContainer, TableRow, Typography, Paper, TableBody, Chip, TableHead, Button, TextField, Snackbar } from "@mui/material";
import AddSprinklerSpecialMode from "./AddSprinklerSpecialMode";
import CloseIcon from '@mui/icons-material/Close';
import LogTable from "../AirConditioner/LogTable";

export class Sprinkler extends Component {

    constructor(props) {
        super(props);

        this.state = {
            specialModes: [],
            startDate: '',
            endDate: '',
            pickedValue: '',
            email: '',
            logData: [],
        };
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleAddSpecialMode = (specialModes) => {
        console.log("add special mode");
        console.log(specialModes);
        this.setState({specialModes: specialModes});
        // TODO: save to database
    }

    handleFormSubmit = () => {
        console.log("form submitted");

    }

    render() {
        
        return (
            <div>
                <Navigation/>
                <img src='/images/arrow.png' alt='arrow' id='arrow' style={{ margin: "55px 0 0 90px", cursor: "pointer" }} onClick={this.handleBackArrow}/>
                <span className='estate-title'>Sprinkler</span>
                <div className='sp-container'>
                    <div id="ac-left-card">
                        <p className='sp-card-title'>Details</p>
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <Typography style={{ fontSize: '1.1em'}}>Off</Typography>
                                <Switch
                                    // checked={item.switchOn}
                                    // onChange={() => this.handleSwitchToggle(item)}
                                />
                            <Typography style={{ fontSize: '1.1em' }}>On</Typography>
                        </div>
                        {/* <p id="special-mode">Add special mode</p> */}
                        <AddSprinklerSpecialMode
                         onAdd={this.handleAddSpecialMode} 
                         isSprinkler='true'/>
                         <TableContainer component={Paper}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Start</TableCell>
                                        <TableCell>End</TableCell>
                                        <TableCell>Day</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {this.state.specialModes.map((item, index) => (
                                        <TableRow key={index}>
                                            <TableCell>{item.start}</TableCell>
                                            <TableCell>{item.end}</TableCell>
                                            <TableCell>
                                                {item.selectedDays.map((day, dayIndex) => (
                                                    <Chip key={dayIndex} label={day} />
                                                ))}
                                            </TableCell>
                                            <TableCell>
                                                <IconButton color="secondary" onClick={() => this.handleDelete(index)}>
                                                    <CloseIcon />
                                                </IconButton>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </div>
                    <div id='sp-right-card'>
                        <p className='sp-card-title'>Switch History</p>
                        <form onSubmit={this.handleFormSubmit} className='sp-container'>
                            <label>
                                Email:
                                <select style={{width: "200px", cursor: "pointer"}}
                                    className="new-real-estate-select"
                                    value={this.state.pickedValue}
                                    onChange={(e) => this.setState({ pickedValue: e.target.value })}>
                                    <option value={this.state.email}>{ this.state.email }</option>
                                    <option value="auto">auto</option>
                                    <option value="none">none</option>
                                </select>
                            </label>
                            <label>
                                Start Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={this.state.startDate} onChange={(e) => this.setState({ startDate: e.target.value })} />
                            </label>
                            <label>
                                End Date:
                                <TextField style={{ backgroundColor: "white" }} type="date" value={this.state.endDate} onChange={(e) => this.setState({ endDate: e.target.value })} />
                            </label>
                            <br />
                            <Button type="submit" id='sp-data-button'>Filter</Button>
                            </form>
                            <LogTable logData={this.state.logData} />
                        </div>
                    </div>
                <Snackbar
                    open={this.state.open}
                    autoHideDuration={2000}
                    onClose={this.handleClose}
                    message={this.state.snackbarMessage}
                    action={this.action}
                />
            </div>
        )
    }
}