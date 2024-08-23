import { defineStore } from "pinia";
import { Email } from "../types/email";
import { getEmails, getEmail } from "../services/emailService";

interface State {
  emails: Email[];
  selectedEmail: Email | null;
  totalPages: number;
  totalElements: number;
  filterTerm: string;
}

export const useEmailStore = defineStore("email", {
  state: (): State => {
    return {
      selectedEmail: null as Email | null,
      emails: [] as Email[],
      totalPages: 0,
      totalElements: 0,
      filterTerm: ""
    };
  },
  actions: {
    async fetchEmail(id: string, filterTerm?: string) {
      this.selectedEmail = await getEmail(id, filterTerm);
    },
    async fetchEmails(params = {}) {
      const response = await getEmails(params);
      this.emails = response.emails || [];
      this.totalPages = response.total_pages || 0;
      this.totalElements = response.total_elements || 0;
    },
  },
});
