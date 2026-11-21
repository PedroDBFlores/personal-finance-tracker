import { createSignal, For, Show } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import { createMovement } from "../api";
import { RouterLink } from "../components/RouterLink";

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
          <div>
            <label class="block text-sm font-medium mb-1" for="movement-description">
              Description
            </label>
            <input
              id="movement-description"
              aria-label="description"
              class="input"
              required
              value={description()}
              onInput={(e) => setDescription(e.currentTarget.value)}
            />
          </div>
          <div class="flex gap-3">
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="movement-amount">
                Amount
              </label>
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
            </div>
            <div class="flex-1">
              <label class="block text-sm font-medium mb-1" for="movement-type">
                Movement type
              </label>
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
            </div>
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
