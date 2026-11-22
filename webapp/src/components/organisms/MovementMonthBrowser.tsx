import { createSignal, For, Show } from "solid-js";
import type { Movement, MovementMonth } from "../../api";
import { listMovementMonths, listMovementsByMonth } from "../../api";
import { Badge } from "../atoms/Badge";
import { EmptyState } from "../atoms/EmptyState";
import { Loading } from "../atoms/Loading";
import { ListRow } from "../molecules/ListRow";

const monthValue = (m: MovementMonth) => `${m.year}-${m.month}`;
const monthLabel = (m: MovementMonth) => `${m.year}-${String(m.month).padStart(2, "0")}`;

export function MovementMonthBrowser(props: { accountId: string }) {
  const [months, setMonths] = createSignal<MovementMonth[] | null>(null);
  const [selected, setSelected] = createSignal<MovementMonth | null>(null);
  const [movements, setMovements] = createSignal<Movement[] | null>(null);
  const [error, setError] = createSignal<string | null>(null);

  async function loadMovements(ym: MovementMonth) {
    setMovements(await listMovementsByMonth(props.accountId, ym.year, ym.month));
  }

  // the component body runs once on mount: load the months and default the
  // selection to the most recent one (the API returns pairs most recent first)
  void (async () => {
    try {
      const list = await listMovementMonths(props.accountId);
      setMonths(list);
      if (list.length > 0) {
        setSelected(list[0]);
        await loadMovements(list[0]);
      } else {
        setMovements([]);
      }
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  })();

  async function onMonthChange(e: Event) {
    const ym = months()?.find((m) => monthValue(m) === (e.currentTarget as HTMLSelectElement).value);
    if (!ym) return;
    setSelected(ym);
    try {
      await loadMovements(ym);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }

  return (
    <div class="flex flex-col gap-2">
      <Show when={error()}>
        <p role="alert" class="text-red-500">
          {error()}
        </p>
      </Show>
      <Show when={months()} fallback={<Loading label="Loading months..." />}>
        {(list) => (
          <Show when={list().length > 0} fallback={<EmptyState>No movements yet</EmptyState>}>
            <label class="block text-sm font-medium">
              Month{" "}
              <select aria-label="month" class="input w-auto" value={monthValue(selected()!)} onChange={onMonthChange}>
                <For each={list()}>
                  {(m) => <option value={monthValue(m)}>{monthLabel(m)}</option>}
                </For>
              </select>
            </label>
            <Show when={movements()} fallback={<Loading label="Loading movements..." />}>
              {(ms) => (
                <ul class="flex flex-col gap-2">
                  <For each={ms()}>
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
              )}
            </Show>
          </Show>
        )}
      </Show>
    </div>
  );
}
