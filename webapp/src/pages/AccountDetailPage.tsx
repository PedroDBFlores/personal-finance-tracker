import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router";
import { Badge, Box, Heading, HStack, List, Text, VStack, Card } from "@chakra-ui/react";
import type { Account, Movement } from "../api";
import { getAccount, listMovements } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function AccountDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [account, setAccount] = useState<Account | null>(null);
  const [movements, setMovements] = useState<Movement[]>([]);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    if (!id) return;
    try {
      const [a, m] = await Promise.all([getAccount(id), listMovements(id)]);
      setAccount(a);
      setMovements(m);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }, [id]);

  useEffect(() => {
    void reload();
  }, [reload]);

  return (
    <Box>
      <HStack justify="space-between" mb={3}>
        <Heading size="lg">{account ? account.name : "Account"}</Heading>
        <RouterLink
          to={`/accounts/${id}/movements/new`}
          px="4"
          py="2"
          borderRadius="md"
          bg="teal.500"
          color="white"
          _hover={{ bg: "teal.600" }}
        >
          Add movement
        </RouterLink>
      </HStack>
      {account && (
        <HStack gap={2} mb={4}>
          <Badge colorPalette="gray">{account.type}</Badge>
          <Badge colorPalette="purple">{account.currency}</Badge>
        </HStack>
      )}
      {error && (
        <Text color="fg.error" role="alert">
          {error}
        </Text>
      )}
      <VStack align="stretch" gap={2}>
        <Heading size="md">Movements</Heading>
        <List.Root gap={2}>
          {movements.map((m) => (
            <List.Item key={m.id}>
              <Card.Root size="sm">
                <Card.Body>
                  <HStack justify="space-between">
                    <Text fontWeight="medium">{m.description}</Text>
                    <HStack gap={2}>
                      <Badge colorPalette={m.type === "Credit" ? "green" : "red"}>{m.type}</Badge>
                      <Text>
                        {m.type === "Debit" ? "-" : "+"}
                        {m.value}
                      </Text>
                    </HStack>
                  </HStack>
                </Card.Body>
              </Card.Root>
            </List.Item>
          ))}
        </List.Root>
        <RouterLink to="/" color="fg.muted" alignSelf="start">
          Back to accounts
        </RouterLink>
      </VStack>
    </Box>
  );
}
