import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  listAccounts,
  getAccount,
  createAccount,
  updateAccount,
  deleteAccount,
  listMovements,
  createMovement,
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

  it("throws ApiError with status and message on failure", async () => {
    mockFetchOnce(404, { error: "resource Account with id 999 not found" });

    await expect(listAccounts()).rejects.toMatchObject({
      status: 404,
      message: "resource Account with id 999 not found",
    });
    await expect(listAccounts()).rejects.toBeInstanceOf(ApiError);
  });
});
