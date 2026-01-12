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
}

export interface AuthResponse {
  token: string;
  user: User;
}
