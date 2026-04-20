// @ts-check

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'devops-reporter',
  tagline: 'Generate static HTML reports from DevOps tool JSON output',
  url: 'https://ndkprd.github.io',
  baseUrl: '/devops-reporter/',
  organizationName: 'ndkprd',
  projectName: 'devops-reporter',
  trailingSlash: false,
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'devops-reporter',
        items: [
          {
            href: 'https://github.com/ndkprd/devops-reporter',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              { label: 'Introduction', to: '/' },
              { label: 'Installation', to: '/installation' },
              { label: 'Custom Templates', to: '/templates' },
            ],
          },
          {
            title: 'Sources',
            items: [
              { label: 'ArgoCD', to: '/sources/argocd' },
              { label: 'Kubeconform', to: '/sources/kubeconform' },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/ndkprd/devops-reporter',
              },
              {
                label: 'Releases',
                href: 'https://github.com/ndkprd/devops-reporter/releases',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} ndkprd. MIT License.`,
      },
    }),
};

module.exports = config;
