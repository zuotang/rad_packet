import { defineStore } from "pinia";
import api from "../api/client";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: localStorage.getItem("token") || "",
    user: null,
  }),
  actions: {
    async login(payload) {
      const inviteCode = localStorage.getItem("invite_code");
      const res = await api.post("/auth/login", payload);
      this.token = res.data.token;
      this.user = res.data.user;
      localStorage.setItem("token", this.token);
      if (inviteCode) {
        try {
          await api.post("/referral/bind", { code: inviteCode });
          localStorage.removeItem("invite_code");
        } catch (_) {
          // keep silent on auto-bind failure, user can manually retry.
        }
      }
    },
    logout() {
      this.token = "";
      this.user = null;
      localStorage.removeItem("token");
    },
  },
});
