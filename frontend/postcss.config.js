const tailwindcss = require('tailwindcss')
const cssnano = require('cssnano')
const autoprefixer = require('autoprefixer')
const purgecss = require('@fullhuman/postcss-purgecss')({
    // Specify the paths to all of the template files in your project 
    content: [
        './index.html'
    ],
    // Include any special characters you're using in this regular expression
    defaultExtractor: content => content.match(/[\w-/:]+(?<!:)/g) || []
})

module.exports = {
    plugins: [
        tailwindcss('./tailwind.js'),
        cssnano({
            preset: 'default',
        }),
        purgecss,
        autoprefixer
    ]
}
