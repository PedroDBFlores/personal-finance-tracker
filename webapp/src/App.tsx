import { useCallback, useEffect, useState } from "react";
import type { Account, Movement } from "./api";
import {
  createAccount,
  createMovement,
  listAccounts,
  listMovements,
} from "./api";

const ACCOUNT_TYPES = ["Bank", "Brokerage", "CreditCard"];
const CURRENCIES = ["USD", "EUR", "GBP"];
const MOVEMENT_TYPES = ["Credit", "Debit"];

export default function App() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [selected, setSelected] = useState<Account | null>(null);

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
    <main style={{ fontFamily: "sans-serif", maxWidth: 720, margin: "2rem auto" }}>
      <h1>Personal Finance Tracker</h1>
      {error && <p role="alert" style={{ color: "red" }}>{error}</p>}
      <AccountForm onCreated={() => void reload()} />
      <h2>Accounts</h2>
      <ul>
        {accounts.map((a) => (
          <li key={a.id}>
            <button
              type="button"
              onClick={() => setSelected(a)}
              style={{ fontWeight: selected?.id === a.id ? "bold" : "normal" }}
            >
              {a.name}
            </button>{" "}
            ({a.type}, {a.currency})
          </li>
        ))}
      </ul>
      {selected && <MovementsPanel account={selected} />}
    </main>
  );
}

function AccountForm({ onCreated }: { onCreated: () => void }) {
  const [name, setName] = useState("");
  const [type, setType] = useState(ACCOUNT_TYPES[0]);
  const [currency, setCurrency] = useState(CURRENCIES[0]);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    try {
      await createAccount({ name, type, currency });
      setName("");
      setError(null);
      onCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSaving(false);
    }
  }

  return (
    <form onSubmit={submit}>
      <h2>New account</h2>
      <label>
        Name{" "}
        <input aria-label="name" value={name} onChange={(e) => setName(e.target.value)} required />
      </label>{" "}
      <label>
        Type{" "}
        <select aria-label="type" value={type} onChange={(e) => setType(e.target.value)}>
          {ACCOUNT_TYPES.map((t) => (
            <option key={t}>{t}</option>
          ))}
        </select>
      </label>{" "}
      <label>
        Currency{" "}
        <select aria-label="currency" value={currency} onChange={(e) => setCurrency(e.target.value)}>
          {CURRENCIES.map((c) => (
            <option key={c}>{c}</option>
          ))}
        </select>
      </label>{" "}
      <button type="submit" disabled={saving}>
        Add account
      </button>
      {error && <p role="alert" style={{ color: "red" }}>{error}</p>}
    </form>
  );
}

function MovementsPanel({ account }: { account: Account }) {
  const [movements, setMovements] = useState<Movement[]>([]);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    try {
      setMovements(await listMovements(account.id));
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }, [account.id]);

  useEffect(() => {
    void reload();
  }, [reload]);

  return (
    <section>
      <h2>Movements of {account.name}</h2>
      {error && <p role="alert" style={{ color: "red" }}>{error}</p>}
      <ul>
        {movements.map((m) => (
          <li key={m.id}>
            {m.description} — {m.value} {m.type}
          </li>
        ))}
      </ul>
      <MovementForm accountId={account.id} onCreated={() => void reload()} />
    </section>
  );
}

function MovementForm({ accountId, onCreated }: { accountId: string; onCreated: () => void }) {
  const [description, setDescription] = useState("");
  const [value, setValue] = useState("");
  const [type, setType] = useState(MOVEMENT_TYPES[0]);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    try {
      await createMovement(accountId, { description, value, type });
      setDescription("");
      setValue("");
      setError(null);
      onCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setSaving(false);
    }
  }

  return (
    <form onSubmit={submit}>
      <label>
        Description{" "}
        <input
          aria-label="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
        />
      </label>{" "}
      <label>
        Amount{" "}
        <input
          aria-label="amount"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder="0.00"
          required
        />
      </label>{" "}
      <label>
        Type{" "}
        <select aria-label="movement type" value={type} onChange={(e) => setType(e.target.value)}>
          {MOVEMENT_TYPES.map((t) => (
            <option key={t}>{t}</option>
          ))}
        </select>
      </label>{" "}
      <button type="submit" disabled={saving}>
        Add movement
      </button>
      {error && <p role="alert" style={{ color: "red" }}>{error}</p>}
    </form>
  );
}
