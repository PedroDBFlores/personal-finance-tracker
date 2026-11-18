import { useCallback, useEffect, useState } from "react";
import { Box, Heading, List, Badge, Text, HStack } from "@chakra-ui/react";
import type { Account } from "../api";
import { listAccounts } from "../api";
import { RouterLink } from "../components/RouterLink";

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    try {
      setAccounts(await listAccounts());
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }, []);

  useEffect(() => {
    void reload();
  }, [reload]);

  return (
    <Box>
      <HStack justify="space-between" mb={3}>
        <Heading size="lg">Accounts</Heading>
        <RouterLink
          to="/accounts/new"
          px="4"
          py="2"
          borderRadius="md"
          bg="teal.500"
          color="white"
          _hover={{ bg: "teal.600" }}
        >
          New account
        </RouterLink>
      </HStack>
      {error && (
        <Text color="fg.error" role="alert">
          {error}
        </Text>
      )}
      <List.Root gap={2}>
        {accounts.map((a) => (
          <List.Item key={a.id}>
            <RouterLink to={`/accounts/${a.id}`} fontWeight="medium" color="teal.500">
              {a.name}
            </RouterLink>{" "}
            <Badge colorPalette="gray">{a.type}</Badge> <Badge colorPalette="purple">{a.currency}</Badge>
          </List.Item>
        ))}
      </List.Root>
    </Box>
  );
}
