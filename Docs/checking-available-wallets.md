
```javascript

console.log('ethereum:', typeof window.ethereum !== 'undefined');
console.log('trustwallet:', typeof window.trustwallet !== 'undefined');
console.log('web3:', typeof window.web3 !== 'undefined');

const allWindowProperties = Object.getOwnPropertyNames(window);
console.log('app props window:', allWindowProperties);

const walletRelatedProperties = allWindowProperties.filter(prop => 
  prop.toLowerCase().includes('wallet') || 
  prop.toLowerCase().includes('web3') || 
  prop.toLowerCase().includes('eth') ||
  prop.toLowerCase().includes('chain')
);
console.log('Props about wallets:', walletRelatedProperties);
```
