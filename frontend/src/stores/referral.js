import { defineStore } from "pinia";
import api from "../api/client";

export const useReferralStore = defineStore("referral", {
  state: () => ({
    myCode: "",
    inviteCount: 0,
    directInvites: [],
  }),
  actions: {
    async fetchStatus() {
      const res = await api.get("/referral/status");
      this.myCode = res.data.my_code;
      this.inviteCount = res.data.invite_count;
      this.directInvites = res.data.direct_invites || [];
    },
    async bind(code) {
      return api.post("/referral/bind", { code });
    },
  },
});
