// wagmi.js
import { http, createConfig } from '@wagmi/vue'
import { sepolia } from '@wagmi/vue/chains'
import { injected } from '@wagmi/vue/connectors'

export const config = createConfig({
    chains: [sepolia],
    connectors: [injected()],
    transports: {
        [sepolia.id]: http(),
    },
})