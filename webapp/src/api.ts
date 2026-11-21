const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export interface Account {
  id: string;
  name: string;
  type: string;
  currency: string;
  created_at: string;
}

export interface CreateAccountInput {
  name: string;
  type: string;
  currency: string;
}

export type UpdateAccountInput = Partial<CreateAccountInput>;

export interface Movement {
  id: string;
  account_id: string;
  description: string;
  value: string; // decimal string, e.g. "25.99"
  type: string;
  date: string;
}

export interface CreateMovementInput {
  description: string;
  value: string;
  type: string;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    let message = `request failed with status ${res.status}`;
    try {
      const body = await res.json();
      if (body?.error) message = body.error;
    } catch {
      // keep default message
    }
    throw new ApiError(res.status, message);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export function listAccounts(): Promise<Account[]> {
  return request<Account[]>("/accounts");
}

export function getAccount(id: string): Promise<Account> {
  return request<Account>(`/accounts/${id}`);
}

export interface AccountBalance {
  account_id: string;
  currency: string;
  movements_balance: string; // decimal string, e.g. "974.01"
  investments_balance: string; // decimal string, e.g. "1000.00"
  total_balance: string; // decimal string, e.g. "1974.01"
}

export function getAccountBalance(accountId: string): Promise<AccountBalance> {
  return request<AccountBalance>(`/accounts/${accountId}/balance`);
}

export function createAccount(input: CreateAccountInput): Promise<Account> {
  return request<Account>("/accounts", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateAccount(id: string, input: UpdateAccountInput): Promise<void> {
  return request<void>(`/accounts/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

export function deleteAccount(id: string): Promise<void> {
  return request<void>(`/accounts/${id}`, { method: "DELETE" });
}

export function listMovements(accountId: string): Promise<Movement[]> {
  return request<Movement[]>(`/accounts/${accountId}/movements`);
}

export function createMovement(accountId: string, input: CreateMovementInput): Promise<Movement> {
  return request<Movement>(`/accounts/${accountId}/movements`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export interface Investment {
  id: string;
  account_id: string;
  name: string;
  type: string;
  current_value: string; // decimal string, e.g. "760.00"
  currency: string;
  purchase_date: string;
  maturity_date?: string | null;
}

export interface CreateInvestmentInput {
  name: string;
  type: string;
  currency: string;
}

export interface InvestmentTransaction {
  id: string;
  investment_id: string;
  type: string;
  amount: string; // decimal string, e.g. "500.00"
  quantity: string; // decimal string, e.g. "10"
  price: string; // decimal string, e.g. "100.00"
  date: string;
  description: string;
}

export interface CreateInvestmentTransactionInput {
  type: string;
  amount: string;
  quantity: string;
  price: string;
  description: string;
}

export function listInvestments(accountId: string): Promise<Investment[]> {
  return request<Investment[]>(`/accounts/${accountId}/investments`);
}

export function getInvestment(accountId: string, id: string): Promise<Investment> {
  return request<Investment>(`/accounts/${accountId}/investments/${id}`);
}

export function createInvestment(accountId: string, input: CreateInvestmentInput): Promise<Investment> {
  return request<Investment>(`/accounts/${accountId}/investments`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function listInvestmentTransactions(accountId: string, investmentId: string): Promise<InvestmentTransaction[]> {
  return request<InvestmentTransaction[]>(`/accounts/${accountId}/investments/${investmentId}/transactions`);
}

export function createInvestmentTransaction(
  accountId: string,
  investmentId: string,
  input: CreateInvestmentTransactionInput,
): Promise<InvestmentTransaction> {
  return request<InvestmentTransaction>(`/accounts/${accountId}/investments/${investmentId}/transactions`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}
