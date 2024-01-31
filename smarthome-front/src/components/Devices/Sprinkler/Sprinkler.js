import React, { Component } from "react";
import './Sprinkler.css';
import { Navigation } from "../../Navigation/Navigation";
import { IconButton, Switch, Table, TableCell, TableContainer, TableRow, Typography, Paper, TableBody, Chip, TableHead } from "@mui/material";
import AddSprinklerSpecialMode from "./AddSprinklerSpecialMode";
import CloseIcon from '@mui/icons-material/Close';

export class Sprinkler extends Component {

    constructor(props) {
        super(props);

        this.state = {
            specialModes: []
        };
    }

    handleBackArrow() {
        window.location.assign("/devices")
    }

    handleAddSpecialMode = (specialModes) => {
        console.log("add special mode");
        console.log(specialModes);
        this.setState({specialModes: specialModes});
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
                </div>
            </div>
        )
    }
}