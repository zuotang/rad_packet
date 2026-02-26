import { defineStore } from "pinia";
import api from "../api/client";

export const useRewardStore = defineStore("reward", {
  state: () => ({
    summary: {
      pending: 0,
      unlocked: 0,
      expired: 0,
    },
  }),
  actions: {
    async fetchSummary() {
      const res = await api.get("/reward/summary");
      this.summary = res.data;
    },
    async unlockPending() {
      const res = await api.post("/reward/unlock");
      if (res?.data?.summary) {
        this.summary = res.data.summary;
      }
      return res;
    },
  },
});
