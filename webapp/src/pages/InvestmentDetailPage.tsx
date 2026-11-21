import { createSignal, For, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import type { Investment, InvestmentTransaction } from "../api";
import { getInvestment, listInvestmentTransactions } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function InvestmentDetailPage() {
  const params = useParams<{ id: string; investmentId: string }>();
  const [investment, setInvestment] = createSignal<Investment | null>(null);
  const [transactions, setTransactions] = createSignal<InvestmentTransaction[]>([]);
  const [error, setError] = createSignal<string | null>(null);

  // the component body runs once on mount
  void (async () => {
    const { id: accountId, investmentId } = params;
    if (!accountId || !investmentId) return;
    try {
      const [inv, txs] = await Promise.all([
        getInvestment(accountId, investmentId),
        listInvestmentTransactions(accountId, investmentId),
      ]);
      setInvestment(inv);
      setTransactions(txs);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  })();

  return (
    <div>
      <div class="flex justify-between items-center mb-3">
        <h2 class="text-xl font-semibold">{investment() ? investment()!.name : "Investment"}</h2>
        <RouterLink
          to={`/accounts/${params.id}/investments/${params.investmentId}/transactions/new`}
          class="btn-primary"
        >
          Add transaction
        </RouterLink>
      </div>
      <Show when={investment()}>
        <div class="flex gap-2 items-center mb-4">
          <span class="badge">{investment()!.type}</span>
          <span class="badge badge-purple">{investment()!.currency}</span>
          <span class="font-medium">Current value: {investment()!.current_value}</span>
        </div>
      </Show>
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <div class="flex flex-col gap-2">
        <h3 class="text-lg font-semibold">Transactions</h3>
        <ul class="flex flex-col gap-2">
          <For each={transactions()}>
            {(tx) => (
              <li>
                <div class="card flex justify-between items-center">
                  <span class="font-medium">{tx.description || tx.type}</span>
                  <span class="flex gap-2 items-center">
                    <span class="badge">{tx.type}</span>
                    <span>
                      {tx.type === "Deposit" || tx.type === "Buy" ? "+" : tx.type === "Withdraw" || tx.type === "Sell" ? "-" : ""}
                      {tx.type === "Deposit" || tx.type === "Withdraw" ? tx.amount : `${tx.quantity} @ ${tx.price}`}
                    </span>
                  </span>
                </div>
              </li>
            )}
          </For>
        </ul>
        <RouterLink to={`/accounts/${params.id}`} class="link-muted self-start">
          Back to account
        </RouterLink>
      </div>
    </div>
  );
}
