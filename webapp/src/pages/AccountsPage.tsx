import { createSignal, For, Show } from "solid-js";
import type { Account } from "../api";
import { listAccounts } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function AccountsPage() {
  const [accounts, setAccounts] = createSignal<Account[]>([]);
  const [error, setError] = createSignal<string | null>(null);

  // the component body runs once on mount
  void (async () => {
    try {
      setAccounts(await listAccounts());
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  })();

  return (
    <div>
      <div class="flex justify-between items-center mb-3">
        <h2 class="text-xl font-semibold">Accounts</h2>
        <RouterLink to="/accounts/new" class="btn-primary">
          New account
        </RouterLink>
      </div>
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <ul class="flex flex-col gap-2">
        <For each={accounts()}>
          {(a) => (
            <li>
              <RouterLink to={`/accounts/${a.id}`} class="link font-medium">
                {a.name}
              </RouterLink>{" "}
              <span class="badge">{a.type}</span> <span class="badge badge-purple">{a.currency}</span>
            </li>
          )}
        </For>
      </ul>
    </div>
  );
}
