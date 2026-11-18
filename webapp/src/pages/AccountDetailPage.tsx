import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router";
import { Badge, Box, Heading, HStack, List, Text, VStack, Card } from "@chakra-ui/react";
import type { Account, Investment, Movement } from "../api";
import { getAccount, listInvestments, listMovements } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function AccountDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [account, setAccount] = useState<Account | null>(null);
  const [movements, setMovements] = useState<Movement[]>([]);
  const [investments, setInvestments] = useState<Investment[]>([]);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    if (!id) return;
    try {
      const [a, m, invs] = await Promise.all([getAccount(id), listMovements(id), listInvestments(id)]);
      setAccount(a);
      setMovements(m);
      setInvestments(invs);
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
      <VStack align="stretch" gap={2} mt={6}>
        <HStack justify="space-between">
          <Heading size="md">Investments</Heading>
          <RouterLink
            to={`/accounts/${id}/investments/new`}
            px="4"
            py="2"
            borderRadius="md"
            bg="teal.500"
            color="white"
            _hover={{ bg: "teal.600" }}
          >
            New investment
          </RouterLink>
        </HStack>
        <List.Root gap={2}>
          {investments.map((inv) => (
            <List.Item key={inv.id}>
              <Card.Root size="sm">
                <Card.Body>
                  <HStack justify="space-between">
                    <RouterLink to={`/accounts/${id}/investments/${inv.id}`} fontWeight="medium" color="teal.500">
                      {inv.name}
                    </RouterLink>
                    <HStack gap={2}>
                      <Badge colorPalette="gray">{inv.type}</Badge>
                      <Text>{inv.current_value}</Text>
                    </HStack>
                  </HStack>
                </Card.Body>
              </Card.Root>
            </List.Item>
          ))}
        </List.Root>
      </VStack>
    </Box>
  );
}
