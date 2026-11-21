import { createSignal, For, Show } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import { createInvestmentTransaction } from "../api";
import { RouterLink } from "../components/atoms/RouterLink";
import { FormField } from "../components/molecules/FormField";

const TRANSACTION_TYPES = ["Buy", "Sell", "Deposit", "Withdraw"];

export default function NewInvestmentTransactionPage() {
  const params = useParams<{ id: string; investmentId: string }>();
  const navigate = useNavigate();
  const [type, setType] = createSignal(TRANSACTION_TYPES[0]);
  const [amount, setAmount] = createSignal("");
  const [quantity, setQuantity] = createSignal("");
  const [price, setPrice] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [error, setError] = createSignal<string | null>(null);
  const [saving, setSaving] = createSignal(false);

  // Buy/Sell trade shares (amount = quantity * price); Deposit/Withdraw move money directly.
  const isTrade = () => type() === "Buy" || type() === "Sell";

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    const { id: accountId, investmentId } = params;
    if (!accountId || !investmentId) return;
    setSaving(true);
    try {
      const computedAmount = isTrade()
        ? ((Number(quantity()) || 0) * (Number(price()) || 0)).toFixed(2)
        : amount();
      await createInvestmentTransaction(accountId, investmentId, {
        type: type(),
        amount: computedAmount,
        quantity: isTrade() ? quantity() : "",
        price: isTrade() ? price() : "",
        description: description(),
      });
      navigate(`/accounts/${accountId}/investments/${investmentId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <div class="card">
      <h2 class="text-xl font-semibold mb-4">New transaction</h2>
      <form onSubmit={submit}>
        <div class="flex flex-col gap-3">
          <FormField id="tx-type" label="Transaction type">
            <select
              id="tx-type"
              aria-label="transaction type"
              class="input"
              value={type()}
              onChange={(e) => setType(e.currentTarget.value)}
            >
              <For each={TRANSACTION_TYPES}>
                {(t) => <option value={t}>{t}</option>}
              </For>
            </select>
          </FormField>
          <Show
            when={isTrade()}
            fallback={
              <FormField id="tx-amount" label="Amount">
                <input
                  id="tx-amount"
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
            }
          >
            <div class="flex gap-3">
              <FormField id="tx-quantity" label="Quantity" class="flex-1">
                <input
                  id="tx-quantity"
                  aria-label="quantity"
                  class="input"
                  type="number"
                  step="0.01"
                  min="0"
                  required
                  value={quantity()}
                  onInput={(e) => setQuantity(e.currentTarget.value)}
                />
              </FormField>
              <FormField id="tx-price" label="Price" class="flex-1">
                <input
                  id="tx-price"
                  aria-label="price"
                  class="input"
                  type="number"
                  step="0.01"
                  min="0"
                  required
                  value={price()}
                  onInput={(e) => setPrice(e.currentTarget.value)}
                />
              </FormField>
            </div>
          </Show>
          <FormField id="tx-description" label="Description">
            <input
              id="tx-description"
              aria-label="description"
              class="input"
              value={description()}
              onInput={(e) => setDescription(e.currentTarget.value)}
            />
          </FormField>
          <Show when={error()}>
            <p role="alert" class="text-red-500">
              {error()}
            </p>
          </Show>
          <div class="flex gap-3">
            <button type="submit" class="btn-primary" disabled={saving()}>
              Add transaction
            </button>
            <RouterLink to={`/accounts/${params.id}/investments/${params.investmentId}`} class="link-muted self-center">
              Cancel
            </RouterLink>
          </div>
        </div>
      </form>
    </div>
  );
}
