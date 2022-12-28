/** @type {import('tailwindcss').Config} */
module.exports = {
    darkMode: 'class',
    content: ['*.html'],
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
    ],
}