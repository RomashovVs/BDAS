import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const HOST = 'localhost';
const PORT = 3000;

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: HOST,
    port: PORT,
  },
})
