import axios from 'axios';
import { Email, EmailResponse } from '../types/email';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

interface GetEmailsParams {
    page?: number;
    size?: number;
    filter?: string;
    sort?: string;
    order?: 'asc' | 'desc';
}

export const getEmail = async (id: string): Promise<Email> => {
    try {
        const response = await axios.get(`${API_BASE_URL}/emails/${id}`);
        return response.data;
    } catch (error) {
        console.error("Error fetching email:", error);
        throw error;
    }
}

export const getEmails = async (params: GetEmailsParams): Promise<EmailResponse> => {
    try {
        const response = await axios.get(`${API_BASE_URL}/emails`, { params });
        console.log("response.data", response.data);
        return response.data;
    } catch (error) {
        console.error("Error fetching emails:", error);
        throw error;
    }
}