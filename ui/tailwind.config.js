module.exports = {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      screens: {
        xxs: '300px',
        xs: '480px',
      },
    },
  },
  important: true,
  corePlugins: {
    preflight: false,
  },
};
