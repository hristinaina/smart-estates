import React from 'react';
import Select from 'react-select';

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
      const options = [
        { value: 'apple', label: 'Apple' },
        { value: 'banana', label: 'Banana' },
        { value: 'orange', label: 'Orange' },
        // Add more options as needed
      ];
      this.setState({ options: options });
    }
    if (type == "rs") {
      const options = [
        { value: 'rs1', label: 'rs1' },
        { value: 'rs2', label: 'rs2' },
        { value: 'rs3', label: 'rs3' },
        // Add more options as needed
      ];
      this.setState({ options: options });
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
