import { createSignal, For, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { createAccount } from "../api";
import { RouterLink } from "../components/atoms/RouterLink";
import { FormField } from "../components/molecules/FormField";

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
          <FormField id="account-name" label="Name">
            <input
              id="account-name"
              aria-label="name"
              class="input"
              required
              value={name()}
              onInput={(e) => setName(e.currentTarget.value)}
            />
          </FormField>
          <div class="flex gap-3">
            <FormField id="account-type" label="Type" class="flex-1">
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
            </FormField>
            <FormField id="account-currency" label="Currency" class="flex-1">
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
            </FormField>
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
