import React from 'react';
import Select from 'react-select';
import RealEstateService from '../../services/RealEstateService';

class SearchSelect extends React.Component {
  state = {
    selectedOptions: [],
    options: []
  };

  componentDidMount() {
    this.updateOptions();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.options !== this.props.options) {
      this.updateOptions();
    }
  }

  updateOptions = async () => {
    const type = this.props.options;
    if (type == "city") {
      const estates = await RealEstateService.getCities();
      const options = []
      for (let i = 0; i < estates.length; i++) {
        const e = estates[i];
        options.push({ value: e, label: e });
      }
      this.setState({ options: options, selectedOptions: [] });
    }
    if (type == "rs") {
      const estates = await RealEstateService.get();
      const options = []
      for (let i = 0; i < estates.length; i++) {
        const e = estates[i];
        options.push({ value: e.Name, label: e.Name });
      }
      this.setState({ options: options, selectedOptions: [] });
    }
  }

  handleChange = (selectedOptions) => {
    this.setState({ selectedOptions });
  };

  render() {
    const { selectedOptions, options } = this.state;

    return (
      <Select 
      id='c-select'
        isMulti
        options={options}
        value={selectedOptions}
        onChange={this.handleChange}
        placeholder="Search and select..."
        menuPortalTarget={document.body} // Attach the menu portal to the body
      />
    );
  }
}

export default SearchSelect;
