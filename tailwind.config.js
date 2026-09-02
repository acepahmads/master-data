/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb', // Primary engineering blue
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554',
        },
        navy: {
          700: '#334155',
          800: '#1e293b',
          850: '#172033',
          900: '#0f172a',
          950: '#0b1120', // Ultra-deep command center navy
        }
      },
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Arial', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'Menlo', 'Monaco', 'Consolas', 'monospace'],
      },
      boxShadow: {
        'xs': '0 1px 2px 0 rgba(0, 0, 0, 0.04)',
        'subtle': '0 1px 3px 0 rgba(15, 23, 42, 0.05), 0 1px 2px 0 rgba(15, 23, 42, 0.03)',
        'card': '0 1px 3px 0 rgba(15, 23, 42, 0.06), 0 1px 2px 0 rgba(15, 23, 42, 0.04)',
        'card-hover': '0 10px 25px -5px rgba(15, 23, 42, 0.08), 0 4px 10px -2px rgba(15, 23, 42, 0.03)',
        'dropdown': '0 10px 25px -5px rgba(15, 23, 42, 0.12), 0 4px 10px -2px rgba(15, 23, 42, 0.06)',
        'modal': '0 25px 50px -12px rgba(15, 23, 42, 0.25)',
        'glow-brand': '0 0 15px -3px rgba(37, 99, 235, 0.35)',
        'glow-emerald': '0 0 15px -3px rgba(16, 185, 129, 0.35)',
      },
      fontSize: {
        '3xs': '0.625rem',  // 10px
        '2xs': '0.6875rem', // 11px
        'xs': '0.75rem',    // 12px
        'sm': '0.8125rem',  // 13px
        'base': '0.875rem', // 14px (standard dense enterprise base)
        'md': '0.9375rem',  // 15px
        'lg': '1rem',       // 16px
        'xl': '1.125rem',   // 18px
        '2xl': '1.25rem',   // 20px
        '3xl': '1.5rem',    // 24px
      },
      borderRadius: {
        'card': '0.625rem', // 10px
      }
    },
  },
  plugins: [],
}
