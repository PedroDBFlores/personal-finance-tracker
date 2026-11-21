import { For, Show } from "solid-js";
import type { Account, AccountBalance } from "../../api";
import { RouterLink } from "../atoms/RouterLink";
import { Badge } from "../atoms/Badge";
import { EmptyState } from "../atoms/EmptyState";
import { ListRow } from "../molecules/ListRow";

type AccountListProps = {
  accounts: Account[];
  balances: Record<string, AccountBalance>;
};

export function AccountList(props: AccountListProps) {
  return (
    <Show when={props.accounts.length > 0} fallback={<EmptyState>No accounts yet</EmptyState>}>
      <ul class="flex flex-col gap-2">
        <For each={props.accounts}>
          {(a) => {
            const balance = () => props.balances[a.id];
            return (
              <ListRow
                left={
                  <RouterLink to={`/accounts/${a.id}`} class="link font-medium">
                    {a.name}
                  </RouterLink>
                }
                right={
                  <>
                    <Badge>{a.type}</Badge>
                    <Badge tone="purple">{a.currency}</Badge>
                    <Show when={balance()}>
                      {(b) => (
                        <span class="font-medium">
                          {b().total_balance} {b().currency}
                        </span>
                      )}
                    </Show>
                  </>
                }
              />
            );
          }}
        </For>
      </ul>
    </Show>
  );
}
