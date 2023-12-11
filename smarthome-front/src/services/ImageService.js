import axios from "axios";

class ImageService {

    static async uploadImage(formData, fileName) {
        try {
            await axios.post('http://localhost:8081/api/upload/' + fileName, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            console.log('Image uploaded successfully');
        } catch (error) {
            console.error('Error uploading image', error);
        }
    }

    static async getImage(fileName) {
        try {
            const result = await axios.get('http://localhost:8081/api/upload/' + fileName);

            console.log('Successful getImage!' + fileName);
            console.log(result);
            console.log(result.data.imageUrl);
            return result.data.imageUrl;
        } catch (error) {
            console.error('Error uploading image', error);
        }
    }
}

export default ImageService;