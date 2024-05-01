/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  safelist:[
    // "border-info",
    // "border-b-2",
    // "border-white",
    // "border-2",
    // "bg-base-content",
    // "border-base-content",
    // "text-black",
    // "px-18"
  ],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["night"]
  }
}

