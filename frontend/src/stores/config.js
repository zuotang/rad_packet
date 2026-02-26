import { defineStore } from "pinia";
import api from "../api/client";

export const useConfigStore = defineStore("config", {
  state: () => ({
    tasks: [],
    rewardTiers: [],
    rawConfigs: {},
    withdrawMin: 60,
  }),
  actions: {
    async fetchBootstrap() {
      const res = await api.get("/config/bootstrap");
      const rawTasks = res.data.tasks || [];
      this.tasks = rawTasks.map((t) => ({
        id: t.id ?? t.ID,
        type: t.type ?? t.Type,
        name: t.name ?? t.Name,
        reward_amount: Number(t.reward_amount ?? t.RewardAmount ?? 0),
      }));
      this.rawConfigs = res.data.configs || {};
      if (this.rawConfigs.withdraw_min) {
        const parsed = Number(this.rawConfigs.withdraw_min);
        if (!Number.isNaN(parsed) && parsed > 0) {
          this.withdrawMin = parsed;
        }
      }
      try {
        this.rewardTiers = JSON.parse(res.data.reward_tiers || "[]");
      } catch {
        this.rewardTiers = [];
      }
    },
  },
});
