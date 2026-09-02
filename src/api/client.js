import axios from 'axios';

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Attach JWT Bearer Token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('iot_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response Interceptor: Handle Data Unwrapping and Global Errors
apiClient.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    const status = error.response ? error.response.status : null;
    const message =
      error.response && error.response.data && error.response.data.message
        ? error.response.data.message
        : error.message || 'An unexpected network error occurred';

    if (status === 401) {
      console.warn('[API Client] Unauthorized request (401) - session expired or invalid token');
      localStorage.removeItem('iot_token');
      sessionStorage.removeItem('iot_token');

      // Prevent redirect loops if already on /login
      if (typeof window !== 'undefined' && window.location && window.location.pathname !== '/login') {
        const currentPath = window.location.pathname + window.location.search;
        const redirectQuery = currentPath && currentPath !== '/' ? `?redirect=${encodeURIComponent(currentPath)}` : '';
        window.location.href = `/login${redirectQuery}`;
      }
    }

    return Promise.reject(new Error(message));
  }
);

export default apiClient;
