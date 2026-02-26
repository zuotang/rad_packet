import { defineStore } from "pinia";
import api from "../api/client";

export const useWalletStore = defineStore("wallet", {
  state: () => ({
    balance: 0,
    frozen: 0,
    ledgers: [],
  }),
  actions: {
    async fetchWallet() {
      const res = await api.get("/wallet?page=1&size=20");
      this.balance = res.data.balance;
      this.frozen = res.data.frozen;
      const rawLedgers = res.data.ledgers || [];
      this.ledgers = rawLedgers.map((l) => ({
        id: l.id ?? l.ID,
        amount: Number(l.amount ?? l.Amount ?? 0),
        type: l.type ?? l.Type,
        ref_type: l.ref_type ?? l.RefType,
        ref_id: l.ref_id ?? l.RefID,
        created_at: l.created_at ?? l.CreatedAt,
      }));
    },
    async applyWithdraw(amount) {
      return api.post("/withdraw/apply", { amount });
    },
  },
});
