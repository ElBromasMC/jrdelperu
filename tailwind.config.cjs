const defaultTheme = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './view/**/*.templ',
  ],
  theme: {
    extend: {
      fontFamily: {
        mono: ['Courier Prime', 'monospace'],
        sans: ['"Fira Code"', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        navy: '#1d2747',
        azure: '#39a0ed',
        chalky: '#efefef',
        livid: '#4c6085',
        darkblue: "#0f151b",
        apple: '#ab2422',
      },
    },
  },
  plugins: [],
  corePlugins: { preFlight: true },
};
