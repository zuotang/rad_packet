import axios from "axios";
import { useAdminStore } from "../stores/admin";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || "/api/admin",
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const admin = useAdminStore();
  if (admin.adminKey) {
    config.headers["X-Admin-Key"] = admin.adminKey;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res.data,
  (error) => Promise.reject(error)
);

export default api;
