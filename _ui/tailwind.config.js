import flowbite from 'flowbite/plugin.js';

/** @type {import('tailwindcss').Config} */
export default {
    plugins: [flowbite],
    theme: {
        extend: {},
    },
    content: [
        "./index.html",
        './src/**/*.{html,svelte,js,ts}',
        "./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}"
    ], // for unused CSS
    variants: {
        extend: {},
    },
    darkMode: 'class',
}

