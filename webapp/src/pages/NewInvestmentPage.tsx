import { createSignal, For, Show } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import { createInvestment } from "../api";
import { RouterLink } from "../components/RouterLink";

const INVESTMENT_TYPES = ["ETF", "Savings", "TermDeposit"];
const CURRENCIES = ["USD", "EUR", "GBP"];

export default function NewInvestmentPage() {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [name, setName] = createSignal("");
  const [type, setType] = createSignal(INVESTMENT_TYPES[0]);
  const [currency, setCurrency] = createSignal(CURRENCIES[0]);
  const [error, setError] = createSignal<string | null>(null);
  const [saving, setSaving] = createSignal(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    const id = params.id;
    if (!id) return;
    setSaving(true);
    try {
      await createInvestment(id, { name: name(), type: type(), currency: currency() });
      navigate(`/accounts/${id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">New investment</h2>
      <form onSubmit={submit}>
        <div class="flex flex-col gap-3">
          <div>
            <label class="block text-sm font-medium mb-1" for="investment-name">
              Name
            </label>
            <input
              id="investment-name"
              aria-label="name"
              class="input"
              required
              value={name()}
              onInput={(e) => setName(e.currentTarget.value)}
            />
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="investment-type">
                Type
              </label>
              <select
                id="investment-type"
                aria-label="type"
                class="input"
                value={type()}
                onChange={(e) => setType(e.currentTarget.value)}
              >
                <For each={INVESTMENT_TYPES}>
                  {(t) => <option value={t}>{t}</option>}
                </For>
              </select>
            </div>
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="investment-currency">
                Currency
              </label>
              <select
                id="investment-currency"
                aria-label="currency"
                class="input"
                value={currency()}
                onChange={(e) => setCurrency(e.currentTarget.value)}
              >
                <For each={CURRENCIES}>
                  {(c) => <option value={c}>{c}</option>}
                </For>
              </select>
            </div>
          </div>
          <Show when={error()}>
            <p role="alert" class="text-red-500">
              {error()}
            </p>
          </Show>
          <div class="flex gap-3">
            <button type="submit" class="btn-primary" disabled={saving()}>
              Add investment
            </button>
            <RouterLink to={`/accounts/${params.id}`} class="link-muted self-center">
              Cancel
            </RouterLink>
          </div>
        </div>
      </form>
    </div>
  );
}
