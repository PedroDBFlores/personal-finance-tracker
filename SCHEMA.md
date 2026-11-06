# Personal Finance Tracker Database Schema 

Personal Finance Tracker uses a KV database to store the data. The models themselves are found in the 'models' folder.

| Model | Key format |
| -------- | ------- |
| Account | account:**{accountId}** | 
| Movement | account:**{accountId}**:movement:**{movementId}** | 
| Investment | account:**{accountId}**:investment:**{investmentId}** | 
| Investment Transaction | account:**{accountId}**:investment:**{investmentId}**:transaction:**{transactionId}** | 
| Tag | tag:**{tagId}** | 