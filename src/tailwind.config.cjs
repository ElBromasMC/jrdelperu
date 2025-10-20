const defaultTheme = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './view/**/*.templ',
  ],
  safelist: [
    {
      pattern: /^(btn|card|input|select|modal|drawer|navbar|menu|dropdown|badge|alert|toast|stat)/,
    },
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Roboto', ...defaultTheme.fontFamily.sans],
        header: ['"Fira Code"', ...defaultTheme.fontFamily.sans],
        footer: ['"Fira Code"', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        navy: '#1d2747',
        azure: '#39a0ed',
        chalky: '#efefef',
        livid: '#585b5e',
        cloud: '#b1b1b2',
        darkblue: "#0f151b",
        apple: '#ab2422',
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        admin: {
          "primary": "#3b82f6",
          "secondary": "#8b5cf6",
          "accent": "#06b6d4",
          "neutral": "#1f2937",
          "base-100": "#ffffff",
          "info": "#0ea5e9",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
        },
      },
      "light", // fallback theme for public pages
    ],
  },
  corePlugins: { preFlight: true },
};
