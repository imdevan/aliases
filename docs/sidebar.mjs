import config from './config.mjs';

const apiReference = {
  label: 'API Reference',
  items: [
    { label: 'app', link: '/api/app' },
    { label: 'bookmark', link: '/api/bookmark' },
    { label: 'config', link: '/api/config' },
    { label: 'domain', link: '/api/domain' },
    { label: 'errors', link: '/api/errors' },
    { label: 'flags', link: '/api/flags' },
    { label: 'ui', link: '/api/ui' },
    { label: 'utils', link: '/api/utils' },
    { label: 'workflow', link: '/api/workflow' },
    {
      label: 'Adapters',
      items: [],
    },
  ],
};

const sidebar = [
  { label: 'bookmark', link: '/' },
  { label: 'Install', link: '/install' },
  { label: 'Commands', items: [
    { label: 'bookmark', link: '/commands/bookmark' },
    { label: 'add', link: '/commands/add' },
    { label: 'completion', link: '/commands/completion' },
    {
      label: 'config',
      items: [
        { label: 'config', link: '/commands/config' },
        { label: 'config init', link: '/commands/config-init' },
      ],
    },
    { label: 'delete', link: '/commands/delete' },
    { label: 'edit', link: '/commands/edit' },
    { label: 'list', link: '/commands/list' },

  ] },
  { label: 'Configuration', link: '/configuration' },
];

const isProduction = process.env.NODE_ENV === 'production';
if (!isProduction) {
  sidebar.push(apiReference);
}

export default sidebar;
