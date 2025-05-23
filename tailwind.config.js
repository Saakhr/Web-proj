/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./**/*.html", "./**/*.templ", "./**/*.go", ],
  theme: { extend: {}, },
  plugins: [require("daisyui")],
  daisyui: {
    logs: false,
    //...
  },
};
