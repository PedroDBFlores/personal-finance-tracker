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
import { createInvestment } from "../api";
import { RouterLink } from "../components/RouterLink";

const INVESTMENT_TYPES = ["ETF", "Savings", "TermDeposit"];
const CURRENCIES = ["USD", "EUR", "GBP"];

export default function NewInvestmentPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [type, setType] = useState(INVESTMENT_TYPES[0]);
  const [currency, setCurrency] = useState(CURRENCIES[0]);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!id) return;
    setSaving(true);
    try {
      await createInvestment(id, { name, type, currency });
      navigate(`/accounts/${id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <Card.Root>
      <Card.Header>
        <Card.Title>New investment</Card.Title>
      </Card.Header>
      <Card.Body>
        <form onSubmit={submit}>
          <VStack gap={3} align="stretch">
            <Field.Root required>
              <Field.Label>Name</Field.Label>
              <Input aria-label="name" value={name} onChange={(e) => setName(e.target.value)} />
            </Field.Root>
            <HStack gap={3}>
              <Field.Root>
                <Field.Label>Type</Field.Label>
                <NativeSelect.Root>
                  <NativeSelect.Field
                    aria-label="type"
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                  >
                    {INVESTMENT_TYPES.map((t) => (
                      <option key={t} value={t}>
                        {t}
                      </option>
                    ))}
                  </NativeSelect.Field>
                </NativeSelect.Root>
              </Field.Root>
              <Field.Root>
                <Field.Label>Currency</Field.Label>
                <NativeSelect.Root>
                  <NativeSelect.Field
                    aria-label="currency"
                    value={currency}
                    onChange={(e) => setCurrency(e.target.value)}
                  >
                    {CURRENCIES.map((c) => (
                      <option key={c} value={c}>
                        {c}
                      </option>
                    ))}
                  </NativeSelect.Field>
                </NativeSelect.Root>
              </Field.Root>
            </HStack>
            {error && (
              <Text color="fg.error" role="alert">
                {error}
              </Text>
            )}
            <HStack gap={3}>
              <Button type="submit" colorPalette="teal" loading={saving}>
                Add investment
              </Button>
              <RouterLink to={`/accounts/${id}`} color="fg.muted">
                Cancel
              </RouterLink>
            </HStack>
          </VStack>
        </form>
      </Card.Body>
    </Card.Root>
  );
}
