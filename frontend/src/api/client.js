import axios from "axios";
import { useAuthStore } from "../stores/auth";
import pinia from "../stores/pinia";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || "/api",
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res.data,
  (error) => {
    if (error?.response?.status === 401) {
      const auth = useAuthStore(pinia);
      auth.logout();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
