import { createSignal, For, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { createAccount } from "../api";
import { RouterLink } from "../components/RouterLink";

const ACCOUNT_TYPES = ["Bank", "Brokerage", "CreditCard"];
const CURRENCIES = ["USD", "EUR", "GBP"];

export default function NewAccountPage() {
  const navigate = useNavigate();
  const [name, setName] = createSignal("");
  const [type, setType] = createSignal(ACCOUNT_TYPES[0]);
  const [currency, setCurrency] = createSignal(CURRENCIES[0]);
  const [error, setError] = createSignal<string | null>(null);
  const [saving, setSaving] = createSignal(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    setSaving(true);
    try {
      await createAccount({ name: name(), type: type(), currency: currency() });
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">New account</h2>
      <form onSubmit={submit}>
        <div class="flex flex-col gap-3">
          <div>
            <label class="block text-sm font-medium mb-1" for="account-name">
              Name
            </label>
            <input
              id="account-name"
              aria-label="name"
              class="input"
              required
              value={name()}
              onInput={(e) => setName(e.currentTarget.value)}
            />
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="account-type">
                Type
              </label>
              <select
                id="account-type"
                aria-label="type"
                class="input"
                value={type()}
                onChange={(e) => setType(e.currentTarget.value)}
              >
                <For each={ACCOUNT_TYPES}>
                  {(t) => <option value={t}>{t}</option>}
                </For>
              </select>
            </div>
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="account-currency">
                Currency
              </label>
              <select
                id="account-currency"
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
              Add account
            </button>
            <RouterLink to="/" class="link-muted self-center">
              Cancel
            </RouterLink>
          </div>
        </div>
      </form>
    </div>
  );
}
