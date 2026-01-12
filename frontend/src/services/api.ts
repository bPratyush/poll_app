import axios from 'axios';

// Use environment variable for API URL
// In production (Render), set VITE_API_URL to your backend URL
// Locally, it defaults to localhost:8080
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  signUp: (data: { username: string; email: string; password: string }) =>
    api.post('/auth/signup', data),
  
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  
  me: () => api.get('/auth/me'),
};

// Poll APIs
export const pollAPI = {
  list: () => api.get('/polls'),
  
  get: (id: number) => api.get(`/polls/${id}`),
  
  create: (data: { title: string; description: string; options: string[] }) =>
    api.post('/polls', data),
  
  update: (id: number, data: { title: string; description: string; options: { id?: number; text: string }[] }) =>
    api.put(`/polls/${id}`, data),
  
  delete: (id: number) => api.delete(`/polls/${id}`),
  
  vote: (pollId: number, optionId: number) =>
    api.post(`/polls/${pollId}/vote`, { option_id: optionId }),
  
  getVoters: (optionId: number) =>
    api.get(`/options/${optionId}/voters`),
};

export default api;
