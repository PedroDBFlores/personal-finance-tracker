import { createSignal, Show } from "solid-js";
import type { Account } from "../api";
import { listAccounts } from "../api";
import { Loading } from "../components/atoms/Loading";
import { PageHeader } from "../components/molecules/PageHeader";
import { AccountList } from "../components/organisms/AccountList";

export default function AccountsPage() {
  const [accounts, setAccounts] = createSignal<Account[] | null>(null);
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
      <PageHeader title="Accounts" action={{ to: "/accounts/new", label: "New account" }} />
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <Show when={accounts()} fallback={<Loading label="Loading accounts..." />}>
        {(list) => <AccountList accounts={list()} />}
      </Show>
    </div>
  );
}
