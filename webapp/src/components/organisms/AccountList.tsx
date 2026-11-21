import { For, Show } from "solid-js";
import type { Account } from "../../api";
import { RouterLink } from "../atoms/RouterLink";
import { Badge } from "../atoms/Badge";
import { EmptyState } from "../atoms/EmptyState";
import { ListRow } from "../molecules/ListRow";

export function AccountList(props: { accounts: Account[] }) {
  return (
    <Show when={props.accounts.length > 0} fallback={<EmptyState>No accounts yet</EmptyState>}>
      <ul class="flex flex-col gap-2">
        <For each={props.accounts}>
          {(a) => (
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
                </>
              }
            />
          )}
        </For>
      </ul>
    </Show>
  );
}
