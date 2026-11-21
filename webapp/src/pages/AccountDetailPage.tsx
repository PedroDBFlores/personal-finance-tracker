import { createSignal, For, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import type { Account, Investment, Movement } from "../api";
import { getAccount, listInvestments, listMovements } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function AccountDetailPage() {
  const params = useParams<{ id: string }>();
  const [account, setAccount] = createSignal<Account | null>(null);
  const [movements, setMovements] = createSignal<Movement[]>([]);
  const [investments, setInvestments] = createSignal<Investment[]>([]);
  const [error, setError] = createSignal<string | null>(null);

  // the component body runs once on mount
  void (async () => {
    const id = params.id;
    if (!id) return;
    try {
      const [a, m, invs] = await Promise.all([getAccount(id), listMovements(id), listInvestments(id)]);
      setAccount(a);
      setMovements(m);
      setInvestments(invs);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  })();

  return (
    <div>
      <div class="flex justify-between items-center mb-3">
        <h2 class="text-xl font-semibold">{account() ? account()!.name : "Account"}</h2>
        <RouterLink to={`/accounts/${params.id}/movements/new`} class="btn-primary">
          Add movement
        </RouterLink>
      </div>
      <Show when={account()}>
        <div class="flex gap-2 mb-4">
          <span class="badge">{account()!.type}</span>
          <span class="badge badge-purple">{account()!.currency}</span>
        </div>
      </Show>
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <div class="flex flex-col gap-2">
        <h3 class="text-lg font-semibold">Movements</h3>
        <ul class="flex flex-col gap-2">
          <For each={movements()}>
            {(m) => (
              <li>
                <div class="card flex justify-between items-center">
                  <span class="font-medium">{m.description}</span>
                  <span class="flex gap-2 items-center">
                    <span class={m.type === "Credit" ? "badge badge-green" : "badge badge-red"}>{m.type}</span>
                    <span>
                      {m.type === "Debit" ? "-" : "+"}
                      {m.value}
                    </span>
                  </span>
                </div>
              </li>
            )}
          </For>
        </ul>
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
        <ul class="flex flex-col gap-2">
          <For each={investments()}>
            {(inv) => (
              <li>
                <div class="card flex justify-between items-center">
                  <RouterLink to={`/accounts/${params.id}/investments/${inv.id}`} class="link font-medium">
                    {inv.name}
                  </RouterLink>
                  <span class="flex gap-2 items-center">
                    <span class="badge">{inv.type}</span>
                    <span>{inv.current_value}</span>
                  </span>
                </div>
              </li>
            )}
          </For>
        </ul>
      </div>
    </div>
  );
}
