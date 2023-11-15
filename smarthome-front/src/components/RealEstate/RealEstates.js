
import React,{ Component, useState } from 'react';
import './RealEstates.css';
import Dialog from '../Dialog/Dialog';
import { NewRealEstate } from './NewRealEstate';

import { Navigation } from '../Navigation/Navigation';

export class RealEstates extends Component {
    constructor(props) {
        super(props);

        this.state = {
            showNewRealEstate: false,
            showApproveDialog: false,
            showDiscardDialog: false,
        };
    }

    handleApprove = () => {
        this.setState({showApproveDialog: true, 
                       showDiscardDialog: false,
                      })
    }

    handleDiscard = () => {
        this.setState({showDiscardDialog: true,
                        showApproveDialog: false,
                      })
    }

    handleConfirmDiscard = () => {
        this.setState({showDiscardDialog: false});
        console.log("Request discarted...")
    }

    handleConfirmApprove = () => {
        this.setState({showApproveDialog: false});
        console.log("Approved");
    }

    handleCancel = () => {
        this.setState({showApproveDialog: false,
                       showDiscardDialog: false,
                      });
        console.log("Cancelled...");
    }

    handleAddRealEstateClick = () => {
        window.location.href = '/new-real-estate';
    }

    render() {
        return (
            <div>
                <Navigation />
                {!this.state.showNewRealEstate && (
                <p id="add-real-estate" onClick={this.handleAddRealEstateClick}>
                    <img alt="." src="/images/plus.png" id="plus" />
                    Add Real-Estate
                </p>
                )}

                {this.state.showNewRealEstate ? (
                <NewRealEstate />
                ) : (
                <div id='real-estates-container'>
                    <div className='real-estate-card'>
                        <img alt='real-estate' src='/images/real_estate_example.png' className='real-estate-img' />
                        <div className='real-estate-info'>
                            <p className='real-estate-title'>Villa B dorm</p>
                            <p className='real-estate-text'>Location: Maldives</p>
                            <p className='real-estate-text'>Square Footage: 102 m2</p>
                            <p className='real-estate-text'>Number of Floors: 2</p>
                            <p className='real-estate-text state-color'>Accepted</p>
                        </div>
                    </div>
                </div>
                )}
                 <div id="bottom-bar">
                    <button className='bottom-bar-btn' id='bottom-bar-approve' onClick={this.handleApprove}>APPROVE</button>
                    <button className='bottom-bar-btn' id='bottom-bar-discard' onClick={this.handleDiscard}>DISCARD</button>
                 </div>

                {this.state.showApproveDialog && (
                <Dialog
                    title="Approve Real Estate Request"
                    message="Are you sure you want to approve selected real-estate request?"
                    onConfirm={this.handleConfirmApprove}
                    onCancel={this.handleCancel}
                />
                )}

                {this.state.showDiscardDialog && (
                <Dialog
                    title="Discard Real Estate Request"
                    message="Are you sure you want to discard selected real-estate request?"
                    onConfirm={this.handleConfirmDiscard}
                    onCancel={this.handleCancel}
                />
                )}
            </div>
        )
    }
}