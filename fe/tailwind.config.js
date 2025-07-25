// tailwind.config.js

const { heroui } = require("@heroui/react");

module.exports = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        institucional: "#9d2348",
      },
    },
  },
  darkMode: "class",
  plugins: [heroui()],
};
