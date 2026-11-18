import { useState } from "react";
import { useNavigate, useParams } from "react-router";
import {
  Button,
  Card,
  Field,
  HStack,
  Input,
  NativeSelect,
  Text,
  VStack,
} from "@chakra-ui/react";
import { createInvestmentTransaction } from "../api";
import { RouterLink } from "../components/RouterLink";

const TRANSACTION_TYPES = ["Buy", "Sell", "Deposit", "Withdraw"];

export default function NewInvestmentTransactionPage() {
  const { id: accountId, investmentId } = useParams<{ id: string; investmentId: string }>();
  const navigate = useNavigate();
  const [type, setType] = useState(TRANSACTION_TYPES[0]);
  const [amount, setAmount] = useState("");
  const [quantity, setQuantity] = useState("");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  // Buy/Sell trade shares (amount = quantity * price); Deposit/Withdraw move money directly.
  const isTrade = type === "Buy" || type === "Sell";

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!accountId || !investmentId) return;
    setSaving(true);
    try {
      const computedAmount = isTrade
        ? ((Number(quantity) || 0) * (Number(price) || 0)).toFixed(2)
        : amount;
      await createInvestmentTransaction(accountId, investmentId, {
        type,
        amount: computedAmount,
        quantity: isTrade ? quantity : "",
        price: isTrade ? price : "",
        description,
      });
      navigate(`/accounts/${accountId}/investments/${investmentId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <Card.Root>
      <Card.Header>
        <Card.Title>New transaction</Card.Title>
      </Card.Header>
      <Card.Body>
        <form onSubmit={submit}>
          <VStack gap={3} align="stretch">
            <Field.Root>
              <Field.Label>Transaction type</Field.Label>
              <NativeSelect.Root>
                <NativeSelect.Field
                  aria-label="transaction type"
                  value={type}
                  onChange={(e) => setType(e.target.value)}
                >
                  {TRANSACTION_TYPES.map((t) => (
                    <option key={t} value={t}>
                      {t}
                    </option>
                  ))}
                </NativeSelect.Field>
              </NativeSelect.Root>
            </Field.Root>
            {isTrade ? (
              <HStack gap={3}>
                <Field.Root required>
                  <Field.Label>Quantity</Field.Label>
                  <Input
                    aria-label="quantity"
                    type="number"
                    step="0.01"
                    min="0"
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                  />
                </Field.Root>
                <Field.Root required>
                  <Field.Label>Price</Field.Label>
                  <Input
                    aria-label="price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={price}
                    onChange={(e) => setPrice(e.target.value)}
                  />
                </Field.Root>
              </HStack>
            ) : (
              <Field.Root required>
                <Field.Label>Amount</Field.Label>
                <Input
                  aria-label="amount"
                  type="number"
                  step="0.01"
                  min="0"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
              </Field.Root>
            )}
            <Field.Root>
              <Field.Label>Description</Field.Label>
              <Input
                aria-label="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </Field.Root>
            {error && (
              <Text color="fg.error" role="alert">
                {error}
              </Text>
            )}
            <HStack gap={3}>
              <Button type="submit" colorPalette="teal" loading={saving}>
                Add transaction
              </Button>
              <RouterLink to={`/accounts/${accountId}/investments/${investmentId}`} color="fg.muted">
                Cancel
              </RouterLink>
            </HStack>
          </VStack>
        </form>
      </Card.Body>
    </Card.Root>
  );
}
