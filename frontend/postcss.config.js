const tailwindcss = require('tailwindcss')
const purgecss = require('@fullhuman/postcss-purgecss')
const cssnano = require('cssnano')
const autoprefixer = require('autoprefixer')

module.exports = {
    plugins: [
        tailwindcss('./tailwind.js'),
        cssnano({
            preset: 'default',
        }),
        purgecss({
            content: ['./index.html']
        }),
        autoprefixer
    ]
}
