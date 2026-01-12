export interface User {
  id: number;
  username: string;
  email: string;
}

export interface Option {
  id: number;
  text: string;
  vote_count: number;
}

export interface Poll {
  id: number;
  title: string;
  description: string;
  creator: User;
  options: Option[];
  created_at: string;
  updated_at: string;
  user_voted_option_id?: number;
  poll_edited_after_vote?: boolean;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Notification {
  id: number;
  message: string;
  type: string;
  poll_id?: number;
  read: boolean;
  created_at: string;
}
