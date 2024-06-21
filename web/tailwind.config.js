/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        montserrat: ["Montserrat", "sans-serif"],
        lora: ["Lora", "sans-serif"],
        roboto: ["Roboto", "sans-serif"],
      },
      colors: {
        primary: "#4A90E2",
        secondary: "#50E3C2",
        accent: "#F5A623",
        lightAccent: "rgba(245, 166, 35, 0.7)",
        background: "#F2F2F2",
        text: "#333333",
        mutualText: "#666666",
      },
    },
  },
  plugins: [],
};
