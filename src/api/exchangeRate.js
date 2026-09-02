import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1'
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const exchangeRateApi = {
  getLatestRate(base = 'USD', quote = 'IDR') {
    return api.get(`/exchange-rates/latest?base=${base}&quote=${quote}`).then(res => res.data);
  },

  refreshRate(base = 'USD', quote = 'IDR') {
    return api.post(`/exchange-rates/refresh?base=${base}&quote=${quote}`).then(res => res.data);
  },

  calculatePricing(payload) {
    return api.post('/exchange-rates/calculate-pricing', payload).then(res => res.data);
  }
};
