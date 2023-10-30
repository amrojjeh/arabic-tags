/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/html/base.tmpl",
    "./ui/html/**/*.tmpl",
    "./ui/static/*.js",
  ],
  theme: {
    extend: {
      boxShadow: {
        "key": "0 2px 0 0 rgba(0,0,0,.7)"
      }
    },
  },
  plugins: [],
}

