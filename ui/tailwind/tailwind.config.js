/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/base.go",
    "./ui/**/*.go",
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

