import { createSignal, For, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import type { Account, AccountBalance, Investment, Movement } from "../api";
import { getAccount, getAccountBalance, listInvestments, listMovements } from "../api";
import { RouterLink } from "../components/atoms/RouterLink";
import { Badge } from "../components/atoms/Badge";
import { EmptyState } from "../components/atoms/EmptyState";
import { Loading } from "../components/atoms/Loading";
import { ListRow } from "../components/molecules/ListRow";
import { PageHeader } from "../components/molecules/PageHeader";

export default function AccountDetailPage() {
  const params = useParams<{ id: string }>();
  const [account, setAccount] = createSignal<Account | null>(null);
  const [balance, setBalance] = createSignal<AccountBalance | null>(null);
  const [movements, setMovements] = createSignal<Movement[] | null>(null);
  const [investments, setInvestments] = createSignal<Investment[] | null>(null);
  const [error, setError] = createSignal<string | null>(null);

  // the component body runs once on mount
  void (async () => {
    const id = params.id;
    if (!id) return;
    try {
      const [a, b, m, invs] = await Promise.all([getAccount(id), getAccountBalance(id), listMovements(id), listInvestments(id)]);
      setAccount(a);
      setBalance(b);
      setMovements(m);
      setInvestments(invs);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  })();

  return (
    <div>
      <PageHeader
        title={account() ? account()!.name : "Account"}
        action={{ to: `/accounts/${params.id}/movements/new`, label: "Add movement" }}
      />
      <Show when={account()}>
        <div class="flex gap-2 mb-4">
          <Badge>{account()!.type}</Badge>
          <Badge tone="purple">{account()!.currency}</Badge>
        </div>
      </Show>
      <Show when={balance()}>
        {(b) => (
          <div class="mb-4">
            <p class="text-2xl font-semibold">
              {b().total_balance} {b().currency}
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              Movements {b().movements_balance} · Investments {b().investments_balance}
            </p>
          </div>
        )}
      </Show>
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <div class="flex flex-col gap-2">
        <h3 class="text-lg font-semibold">Movements</h3>
        <Show when={movements()} fallback={<Loading label="Loading movements..." />}>
          {(list) => (
            <Show when={list().length > 0} fallback={<EmptyState>No movements yet</EmptyState>}>
              <ul class="flex flex-col gap-2">
                <For each={list()}>
                  {(m) => (
                    <ListRow
                      left={<span class="font-medium">{m.description}</span>}
                      right={
                        <>
                          <Badge tone={m.type === "Credit" ? "green" : "red"}>{m.type}</Badge>
                          <span>
                            {m.type === "Debit" ? "-" : "+"}
                            {m.value}
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
        <RouterLink to="/" class="link-muted self-start">
          Back to accounts
        </RouterLink>
      </div>
      <div class="flex flex-col gap-2 mt-6">
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-semibold">Investments</h3>
          <RouterLink to={`/accounts/${params.id}/investments/new`} class="btn-primary">
            New investment
          </RouterLink>
        </div>
        <Show when={investments()} fallback={<Loading label="Loading investments..." />}>
          {(list) => (
            <Show when={list().length > 0} fallback={<EmptyState>No investments yet</EmptyState>}>
              <ul class="flex flex-col gap-2">
                <For each={list()}>
                  {(inv) => (
                    <ListRow
                      left={
                        <RouterLink to={`/accounts/${params.id}/investments/${inv.id}`} class="link font-medium">
                          {inv.name}
                        </RouterLink>
                      }
                      right={
                        <>
                          <Badge>{inv.type}</Badge>
                          <span>{inv.current_value}</span>
                        </>
                      }
                    />
                  )}
                </For>
              </ul>
            </Show>
          )}
        </Show>
      </div>
    </div>
  );
}
