// prettier.config.js
module.exports = {
  plugins: ["prettier-plugin-tailwindcss", "prettier-plugin-go-template"],
  overrides: [
    {
      "files": ["*.templ"],
      "options": {
        "parser": "go-template",
      },
    },
  ],
};
