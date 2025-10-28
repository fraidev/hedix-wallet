# Hedix Wallet

A simple wallet application that supports BTC, ETH, and USD transactions with ledger-based transaction tracking.

## Usage

### File Mode

Process transactions from a file:

```bash
go run --file example.txt
```

### Interactive Mode (Default)

Run the application without arguments to enter interactive mode:

```bash
go run
```

You'll be prompted to enter transactions in the following format:

```
<DEPOSIT|WITHDRAW> <BTC|ETH|USD> <amount>
```

**Examples:**

```
> DEPOSIT BTC 1.5
Transaction successful
Current State: BTC: 1.50000000 | ETH: 0.00000000 | USD: 0.00

> DEPOSIT ETH 2.0
Transaction successful
Current State: BTC: 1.50000000 | ETH: 2.00000000 | USD: 0.00

> WITHDRAW BTC 0.5
Transaction successful
Current State: BTC: 1.00000000 | ETH: 2.00000000 | USD: 0.00

> WITHDRAW BTC 5.0
Transaction failed: insufficient funds for withdrawal: requested 5.00, available 1.00
Current State: BTC: 1.00000000 | ETH: 2.00000000 | USD: 0.00
```

Press `Ctrl+C` or `Ctrl+D` to exit.
