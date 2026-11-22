import { createSignal, For, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import type { Investment, InvestmentTransaction } from "../api";
import { getInvestment, listInvestmentTransactions } from "../api";
import { RouterLink } from "../components/atoms/RouterLink";
import { Badge } from "../components/atoms/Badge";
import { EmptyState } from "../components/atoms/EmptyState";
import { Loading } from "../components/atoms/Loading";
import { ListRow } from "../components/molecules/ListRow";
import { PageHeader } from "../components/molecules/PageHeader";

export default function InvestmentDetailPage() {
  const params = useParams<{ id: string; investmentId: string }>();
  const [investment, setInvestment] = createSignal<Investment | null>(null);
  const [transactions, setTransactions] = createSignal<InvestmentTransaction[] | null>(null);
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
      <PageHeader
        title={investment() ? investment()!.name : "Investment"}
        action={{
          to: `/accounts/${params.id}/investments/${params.investmentId}/transactions/new`,
          label: "Add transaction",
        }}
      />
      <Show when={investment()}>
        <div class="flex gap-2 items-center mb-4">
          <Badge>{investment()!.type}</Badge>
          <Badge tone="purple">{investment()!.currency}</Badge>
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
        <Show when={transactions()} fallback={<Loading label="Loading transactions..." />}>
          {(list) => (
            <Show when={list().length > 0} fallback={<EmptyState>No transactions yet</EmptyState>}>
              <ul class="flex flex-col gap-2">
                <For each={list()}>
                  {(tx) => (
                    <ListRow
                      left={<span class="font-medium">{tx.description || tx.type}</span>}
                      right={
                        <>
                          <Badge tone={tx.type === "Buy" || tx.type === "Deposit" ? "green" : "red"}>{tx.type}</Badge>
                          <span>
                            {tx.type === "Deposit" || tx.type === "Buy" ? "+" : tx.type === "Withdraw" || tx.type === "Sell" ? "-" : ""}
                            {tx.type === "Deposit" || tx.type === "Withdraw" ? tx.amount : `${tx.quantity} @ ${tx.price}`}
                          </span>
                        </>
                      }
                    />
                  )}
                </For>
              </ul>
            </Show>
          )}
        </Show>
        <RouterLink to={`/accounts/${params.id}`} class="link-muted self-start">
          Back to account
        </RouterLink>
      </div>
    </div>
  );
}
