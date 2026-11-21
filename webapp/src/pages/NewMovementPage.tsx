import { createSignal, For, Show } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import { createMovement } from "../api";
import { RouterLink } from "../components/atoms/RouterLink";
import { FormField } from "../components/molecules/FormField";

const MOVEMENT_TYPES = ["Credit", "Debit"];

export default function NewMovementPage() {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [description, setDescription] = createSignal("");
  const [amount, setAmount] = createSignal("");
  const [type, setType] = createSignal(MOVEMENT_TYPES[0]);
  const [error, setError] = createSignal<string | null>(null);
  const [saving, setSaving] = createSignal(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    const id = params.id;
    if (!id) return;
    setSaving(true);
    try {
      await createMovement(id, { description: description(), value: amount(), type: type() });
      navigate(`/accounts/${id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">New movement</h2>
      <form onSubmit={submit}>
        <div class="flex flex-col gap-3">
          <FormField id="movement-description" label="Description">
            <input
              id="movement-description"
              aria-label="description"
              class="input"
              required
              value={description()}
              onInput={(e) => setDescription(e.currentTarget.value)}
            />
          </FormField>
          <div class="flex gap-3">
            <FormField id="movement-amount" label="Amount" class="flex-1">
              <input
                id="movement-amount"
                aria-label="amount"
                class="input"
                type="number"
                step="0.01"
                min="0"
                required
                value={amount()}
                onInput={(e) => setAmount(e.currentTarget.value)}
              />
            </FormField>
            <FormField id="movement-type" label="Movement type" class="flex-1">
              <select
                id="movement-type"
                aria-label="movement type"
                class="input"
                value={type()}
                onChange={(e) => setType(e.currentTarget.value)}
              >
                <For each={MOVEMENT_TYPES}>
                  {(t) => <option value={t}>{t}</option>}
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
              Add movement
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
