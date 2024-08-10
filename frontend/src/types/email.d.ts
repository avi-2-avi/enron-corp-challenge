export interface Email {
    id: string;
    date: string;
    from: string;
    to: string;
    subject: string;
    content: string;
    path: string;
    timestamp: string;
  }
  
  export interface EmailResponse {
    page: number;
    size: number;
    total_elements: number;
    total_pages: number;
    emails: Email[];
  }