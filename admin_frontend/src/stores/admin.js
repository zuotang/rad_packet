import { defineStore } from "pinia";

export const useAdminStore = defineStore("admin", {
  state: () => ({
    adminKey: localStorage.getItem("admin_key") || "",
  }),
  actions: {
    setKey(key) {
      this.adminKey = key || "";
      localStorage.setItem("admin_key", this.adminKey);
    },
    logout() {
      this.adminKey = "";
      localStorage.removeItem("admin_key");
    },
  },
});
