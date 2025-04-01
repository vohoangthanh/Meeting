/**
 * auth.js - Handles wallet authentication and blockchain interaction
 */

class Web3Auth {
    constructor() {
        this._isConnected = false;
        this.userAddress = null;
        this.privateKey = null;
        this.web3 = null;
        this.accountsChangedListeners = [];
        
        // Check if there's a persisted connection when instantiating
        this._loadPersistedConnection();
    }

    /**
     * Connect to wallet using private key
     * @param {string} privateKey - The private key to use for connection
     * @returns {Promise<Object>} Connection result
     */
    async connect(privateKey) {
        try {
            if (!privateKey || privateKey.trim() === '') {
                return { success: false, error: 'No private key provided' };
            }

            // Clean up the private key if needed
            if (privateKey.startsWith('0x')) {
                this.privateKey = privateKey;
            } else {
                this.privateKey = '0x' + privateKey;
            }

            // Create a new wallet instance
            const wallet = new ethers.Wallet(this.privateKey);
            this.userAddress = wallet.address;
            
            // Set web3 provider if available
            if (window.ethereum) {
                const provider = new ethers.providers.Web3Provider(window.ethereum);
                this.web3 = provider;
            } else {
                // Use a fallback provider (BSC testnet in this case)
                this.web3 = new ethers.providers.JsonRpcProvider('https://bsc-testnet-rpc.publicnode.com');
            }

            this._isConnected = true;
            
            // Persist connection state
            this._persistConnection();
            
            this.notifyAccountsChanged(this.userAddress);
            
            console.log('Successfully connected wallet:', this.getShortAddress());
            return { success: true, address: this.userAddress };
        } catch (error) {
            console.error('Error connecting wallet:', error);
            this._isConnected = false;
            this.userAddress = null;
            this.privateKey = null;
            return { success: false, error: error.message };
        }
    }

    /**
     * Disconnect the wallet
     */
    disconnect() {
        this._isConnected = false;
        this.userAddress = null;
        this.privateKey = null;
        // Clear persisted state
        localStorage.removeItem('wallet_connected');
        localStorage.removeItem('wallet_address');
        this.notifyAccountsChanged(null);
        console.log('Wallet disconnected');
    }

    /**
     * Get wallet connection status
     * @returns {boolean} Whether wallet is connected
     */
    isConnected() {
        return this._isConnected;
    }

    /**
     * Get shortened address for display
     * @returns {string} Shortened address or empty string if not connected
     */
    getShortAddress() {
        if (!this.userAddress) return '';
        return this.userAddress.substring(0, 6) + '...' + this.userAddress.substring(this.userAddress.length - 4);
    }

    /**
     * Get the full wallet address
     * @returns {string|null} Full wallet address or null if not connected
     */
    getAddress() {
        return this.userAddress;
    }

    /**
     * Add listener for accounts changed event
     * @param {Function} listener - Callback function when accounts change
     */
    addAccountsChangedListener(listener) {
        if (typeof listener === 'function') {
            this.accountsChangedListeners.push(listener);
        }
    }

    /**
     * Remove a listener for accounts changed event
     * @param {Function} listener - The listener to remove
     */
    removeAccountsChangedListener(listener) {
        this.accountsChangedListeners = this.accountsChangedListeners.filter(l => l !== listener);
    }

    /**
     * Notify all listeners of accounts changed
     * @param {string|null} address - The new address or null if disconnected
     * @private
     */
    notifyAccountsChanged(address) {
        for (const listener of this.accountsChangedListeners) {
            try {
                listener(address);
            } catch (error) {
                console.error('Error in accounts changed listener:', error);
            }
        }
    }

    /**
     * Sign a message with the connected wallet
     * @param {string} message - The message to sign
     * @returns {Promise<string>} The signature
     * @throws {Error} If wallet is not connected
     */
    async signMessage(message) {
        if (!this.isConnected || !this.privateKey) {
            throw new Error('Wallet not connected');
        }

        const wallet = new ethers.Wallet(this.privateKey);
        return wallet.signMessage(message);
    }

    /**
     * Sign a transaction
     * @param {Object} tx - The transaction object
     * @returns {Promise<string>} The signed transaction
     * @throws {Error} If wallet is not connected
     */
    async signTransaction(tx) {
        if (!this.isConnected || !this.privateKey) {
            throw new Error('Wallet not connected');
        }

        const wallet = new ethers.Wallet(this.privateKey, this.web3);
        return wallet.signTransaction(tx);
    }
    
    /**
     * Persist the current connection state in localStorage
     * @private
     */
    _persistConnection() {
        if (this._isConnected && this.userAddress) {
            localStorage.setItem('wallet_connected', 'true');
            localStorage.setItem('wallet_address', this.userAddress);
            // We don't store the private key in localStorage for security
            // It's already stored separately by the application
        }
    }
    
    /**
     * Load persisted connection state from localStorage on initialization
     * @private
     */
    _loadPersistedConnection() {
        const connected = localStorage.getItem('wallet_connected') === 'true';
        const address = localStorage.getItem('wallet_address');
        const privateKey = localStorage.getItem('privateKey'); // Using the existing privateKey storage
        
        if (connected && address && privateKey) {
            console.log('Found persisted wallet connection, restoring...');
            // Reconnect using the stored private key
            this.connect(privateKey).then(result => {
                if (result.success) {
                    console.log('Wallet connection restored successfully');
                } else {
                    console.error('Failed to restore wallet connection:', result.error);
                }
            });
        }
    }
}

// Create a singleton instance
const auth = new Web3Auth();

export { auth };
