import axios from "axios";
import { Email, EmailResponse, GetEmailsParams } from "../types/email";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export const getEmail = async (
  id: string,
  filterTerm?: string
): Promise<Email> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/emails/${id}`, {
      params: { filter: filterTerm },
    });
    return response.data;
  } catch (error) {
    console.error("Error fetching email:", error);
    throw error;
  }
};

export const getEmails = async (
  params: GetEmailsParams
): Promise<EmailResponse> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/emails`, { params });
    return response.data;
  } catch (error) {
    console.error("Error fetching emails:", error);
    throw error;
  }
};
