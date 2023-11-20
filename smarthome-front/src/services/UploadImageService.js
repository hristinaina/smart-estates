import axios from "axios";

class UploadImageService {

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
}

export default UploadImageService;