import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  listAccounts,
  getAccount,
  getAccountBalance,
  createAccount,
  updateAccount,
  deleteAccount,
  listMovements,
  createMovement,
  listMovementMonths,
  listMovementsByMonth,
  listInvestments,
  getInvestment,
  createInvestment,
  listInvestmentTransactions,
  createInvestmentTransaction,
  ApiError,
} from "./api";

const base = "http://localhost:8080";

function mockFetchOnce(status: number, body: unknown) {
  const fetchMock = vi.fn().mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  } as Response);
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
}

beforeEach(() => {
  vi.unstubAllGlobals();
});

describe("api client", () => {
  it("lists accounts", async () => {
    const accounts = [{ id: "1", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" }];
    const fetchMock = mockFetchOnce(200, accounts);

    const result = await listAccounts();

    expect(result).toEqual(accounts);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts`, expect.anything());
  });

  it("gets a single account", async () => {
    const account = { id: "1", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" };
    const fetchMock = mockFetchOnce(200, account);

    const result = await getAccount("1");

    expect(result).toEqual(account);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/1`, expect.anything());
  });

  it("gets an account balance", async () => {
    const balance = {
      account_id: "1",
      currency: "EUR",
      movements_balance: "974.01",
      investments_balance: "1000.00",
      total_balance: "1974.01",
    };
    const fetchMock = mockFetchOnce(200, balance);

    const result = await getAccountBalance("1");

    expect(result).toEqual(balance);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/1/balance`, expect.anything());
  });

  it("creates an account with a JSON body", async () => {
    const created = { id: "5", name: "Bank A", type: "Bank", currency: "EUR", created_at: "2024-01-02T03:04:05Z" };
    const fetchMock = mockFetchOnce(201, created);

    const result = await createAccount({ name: "Bank A", type: "Bank", currency: "EUR" });

    expect(result).toEqual(created);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts`);
    expect(init.method).toBe("POST");
    expect(JSON.parse(init.body as string)).toEqual({ name: "Bank A", type: "Bank", currency: "EUR" });
  });

  it("sends only the provided fields on update", async () => {
    const fetchMock = mockFetchOnce(204, null);

    await updateAccount("3", { name: "New Name" });

    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts/3`);
    expect(init.method).toBe("PATCH");
    expect(JSON.parse(init.body as string)).toEqual({ name: "New Name" });
  });

  it("deletes an account", async () => {
    const fetchMock = mockFetchOnce(204, null);

    await deleteAccount("4");

    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts/4`);
    expect(init.method).toBe("DELETE");
  });

  it("lists movements of an account", async () => {
    const movements = [{ id: "1", account_id: "2", description: "Salary", value: "1000.00", type: "Credit", date: "2024-01-02T03:04:05Z" }];
    const fetchMock = mockFetchOnce(200, movements);

    const result = await listMovements("2");

    expect(result).toEqual(movements);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/2/movements`, expect.anything());
  });

  it("creates a movement", async () => {
    const created = { id: "9", account_id: "2", description: "Salary", value: "1000.00", type: "Credit", date: "2024-01-02T03:04:05Z" };
    const fetchMock = mockFetchOnce(201, created);

    const result = await createMovement("2", { description: "Salary", value: "1000.00", type: "Credit" });

    expect(result).toEqual(created);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts/2/movements`);
    expect(init.method).toBe("POST");
  });

  it("lists the movement months of an account", async () => {
    const months = [{ year: 2025, month: 2 }, { year: 2025, month: 1 }];
    const fetchMock = mockFetchOnce(200, months);

    const result = await listMovementMonths("1");

    expect(result).toEqual(months);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/1/movements/months`, expect.anything());
  });

  it("lists the movements of an account for a year/month", async () => {
    const movements = [{ id: "1", account_id: "1", description: "Salary", value: "1000.00", type: "Credit", date: "2025-01-02T03:04:05Z" }];
    const fetchMock = mockFetchOnce(200, movements);

    const result = await listMovementsByMonth("1", 2025, 1);

    expect(result).toEqual(movements);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/1/movements/by-month/2025/1`, expect.anything());
  });

  it("lists investments of an account", async () => {
    const investments = [
      { id: "inv-1", account_id: "2", name: "Vanguard ETF", type: "ETF", current_value: "760.00", currency: "USD", purchase_date: "2024-01-02T03:04:05Z", maturity_date: null },
    ];
    const fetchMock = mockFetchOnce(200, investments);

    const result = await listInvestments("2");

    expect(result).toEqual(investments);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/2/investments`, expect.anything());
  });

  it("gets a single investment", async () => {
    const investment = { id: "inv-1", account_id: "2", name: "Vanguard ETF", type: "ETF", current_value: "760.00", currency: "USD", purchase_date: "2024-01-02T03:04:05Z", maturity_date: null };
    const fetchMock = mockFetchOnce(200, investment);

    const result = await getInvestment("2", "inv-1");

    expect(result).toEqual(investment);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/2/investments/inv-1`, expect.anything());
  });

  it("creates an investment", async () => {
    const created = { id: "inv-9", account_id: "2", name: "Vanguard ETF", type: "ETF", current_value: "0.00", currency: "USD", purchase_date: "2024-01-02T03:04:05Z", maturity_date: null };
    const fetchMock = mockFetchOnce(201, created);

    const result = await createInvestment("2", { name: "Vanguard ETF", type: "ETF", currency: "USD" });

    expect(result).toEqual(created);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts/2/investments`);
    expect(init.method).toBe("POST");
    expect(JSON.parse(init.body as string)).toEqual({ name: "Vanguard ETF", type: "ETF", currency: "USD" });
  });

  it("lists transactions of an investment", async () => {
    const txs = [
      { id: "tx-1", investment_id: "inv-1", type: "Buy", amount: "0.00", quantity: "10.00", price: "100.00", date: "2024-01-02T03:04:05Z", description: "Initial Buy" },
    ];
    const fetchMock = mockFetchOnce(200, txs);

    const result = await listInvestmentTransactions("2", "inv-1");

    expect(result).toEqual(txs);
    expect(fetchMock).toHaveBeenCalledWith(`${base}/accounts/2/investments/inv-1/transactions`, expect.anything());
  });

  it("creates an investment transaction", async () => {
    const created = { id: "tx-9", investment_id: "inv-1", type: "Deposit", amount: "500.00", quantity: "0.00", price: "0.00", date: "2024-01-02T03:04:05Z", description: "Deposit" };
    const fetchMock = mockFetchOnce(201, created);

    const input = { type: "Deposit", amount: "500.00", quantity: "", price: "", description: "Deposit" };
    const result = await createInvestmentTransaction("2", "inv-1", input);

    expect(result).toEqual(created);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe(`${base}/accounts/2/investments/inv-1/transactions`);
    expect(init.method).toBe("POST");
    expect(JSON.parse(init.body as string)).toEqual(input);
  });

  it("throws ApiError with status and message on failure", async () => {
    mockFetchOnce(404, { error: "resource Account with id 999 not found" });

    await expect(listAccounts()).rejects.toMatchObject({
      status: 404,
      message: "resource Account with id 999 not found",
    });
    await expect(listAccounts()).rejects.toBeInstanceOf(ApiError);
  });
});
