class ConsumptionService {

    convertOptionsToStrings(selectedOptions) {
        const options = []
        for (let i = 0; i < selectedOptions.length; i++) {
            const e = selectedOptions[i];
            options.push(String(e.value));
        }
        return options;
    }

    async getConsumptionGraphDataForDropdownSelect(queryType, type, options, time) {
        const selectedOptions = this.convertOptionsToStrings(options);
        //console.log(selectedOptions);
        try {
            const response = await fetch('http://localhost:8081/api/consumption/selected-time', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ type, selectedOptions, time, queryType })
            })
            // console.log(response)

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }

    async getConsumptionGraphDataForDates(queryType, type, options, start, end) {
        const selectedOptions = this.convertOptionsToStrings(options);
        //console.log(selectedOptions);
        try {
            const response = await fetch('http://localhost:8081/api/consumption/selected-date', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ type, selectedOptions, start, end, queryType })
            })
            // console.log(response)

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }

    async getRatioGraphDataForDropdownSelect(type, options, time) {
        const selectedOptions = this.convertOptionsToStrings(options);
        //console.log(selectedOptions);
        try {
            const response = await fetch('http://localhost:8081/api/consumption/ratio/selected-time', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ type, selectedOptions, time })
            })
            // console.log(response)

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }

    async getRatioGraphDataForDates(type, options, start, end) {
        const selectedOptions = this.convertOptionsToStrings(options);
        //console.log(selectedOptions);
        try {
            const response = await fetch('http://localhost:8081/api/consumption/ratio/selected-date', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ type, selectedOptions, start, end })
            })
            // console.log(response)

            if (response.ok) {
                const data = await response.json();
                console.log(data)
                return { result: data };
            } else {
                const data = await response.json();
                return { result: false, error: data.error };
            }
        } catch (error) {
            console.error('Greška :', error);
            return { result: false, error: 'Network error' };
        }
    }
}

export default new ConsumptionService();