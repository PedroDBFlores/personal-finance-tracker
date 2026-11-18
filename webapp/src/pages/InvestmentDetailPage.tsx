import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router";
import { Badge, Box, Heading, HStack, List, Text, VStack, Card } from "@chakra-ui/react";
import type { Investment, InvestmentTransaction } from "../api";
import { getInvestment, listInvestmentTransactions } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function InvestmentDetailPage() {
  const { id: accountId, investmentId } = useParams<{ id: string; investmentId: string }>();
  const [investment, setInvestment] = useState<Investment | null>(null);
  const [transactions, setTransactions] = useState<InvestmentTransaction[]>([]);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    if (!accountId || !investmentId) return;
    try {
      const [inv, txs] = await Promise.all([
        getInvestment(accountId, investmentId),
        listInvestmentTransactions(accountId, investmentId),
      ]);
      setInvestment(inv);
      setTransactions(txs);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }, [accountId, investmentId]);

  useEffect(() => {
    void reload();
  }, [reload]);

  return (
    <Box>
      <HStack justify="space-between" mb={3}>
        <Heading size="lg">{investment ? investment.name : "Investment"}</Heading>
        <RouterLink
          to={`/accounts/${accountId}/investments/${investmentId}/transactions/new`}
          px="4"
          py="2"
          borderRadius="md"
          bg="teal.500"
          color="white"
          _hover={{ bg: "teal.600" }}
        >
          Add transaction
        </RouterLink>
      </HStack>
      {investment && (
        <HStack gap={2} mb={4}>
          <Badge colorPalette="gray">{investment.type}</Badge>
          <Badge colorPalette="purple">{investment.currency}</Badge>
          <Text fontWeight="medium">Current value: {investment.current_value}</Text>
        </HStack>
      )}
      {error && (
        <Text color="fg.error" role="alert">
          {error}
        </Text>
      )}
      <VStack align="stretch" gap={2}>
        <Heading size="md">Transactions</Heading>
        <List.Root gap={2}>
          {transactions.map((tx) => (
            <List.Item key={tx.id}>
              <Card.Root size="sm">
                <Card.Body>
                  <HStack justify="space-between">
                    <Text fontWeight="medium">{tx.description || tx.type}</Text>
                    <HStack gap={2}>
                      <Badge colorPalette="gray">{tx.type}</Badge>
                      <Text>
                        {tx.type === "Deposit" || tx.type === "Buy" ? "+" : tx.type === "Withdraw" || tx.type === "Sell" ? "-" : ""}
                        {tx.type === "Deposit" || tx.type === "Withdraw" ? tx.amount : `${tx.quantity} @ ${tx.price}`}
                      </Text>
                    </HStack>
                  </HStack>
                </Card.Body>
              </Card.Root>
            </List.Item>
          ))}
        </List.Root>
        <RouterLink to={`/accounts/${accountId}`} color="fg.muted" alignSelf="start">
          Back to account
        </RouterLink>
      </VStack>
    </Box>
  );
}
