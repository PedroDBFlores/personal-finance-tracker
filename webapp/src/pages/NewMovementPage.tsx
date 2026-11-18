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
import { createMovement } from "../api";
import { RouterLink } from "../components/RouterLink";

const MOVEMENT_TYPES = ["Credit", "Debit"];

export default function NewMovementPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [description, setDescription] = useState("");
  const [amount, setAmount] = useState("");
  const [type, setType] = useState(MOVEMENT_TYPES[0]);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!id) return;
    setSaving(true);
    try {
      await createMovement(id, { description, value: amount, type });
      navigate(`/accounts/${id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
      setSaving(false);
    }
  }

  return (
    <Card.Root>
      <Card.Header>
        <Card.Title>New movement</Card.Title>
      </Card.Header>
      <Card.Body>
        <form onSubmit={submit}>
          <VStack gap={3} align="stretch">
            <Field.Root required>
              <Field.Label>Description</Field.Label>
              <Input
                aria-label="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </Field.Root>
            <HStack gap={3}>
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
              <Field.Root>
                <Field.Label>Movement type</Field.Label>
                <NativeSelect.Root>
                  <NativeSelect.Field
                    aria-label="movement type"
                    value={type}
                    onChange={(e) => setType(e.target.value)}
                  >
                    {MOVEMENT_TYPES.map((t) => (
                      <option key={t} value={t}>
                        {t}
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
                Add movement
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
