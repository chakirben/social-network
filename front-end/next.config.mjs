/** @type {import('next')} */
const isDev = process.env.NODE_ENV === 'development';

const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: isDev
          ? 'http://127.0.0.1:8080/api/:path*'
          : 'http://backend:8080/api/:path*',
      },
      {
        source: '/uploads/:path*',
        destination: isDev
          ? 'http://127.0.0.1:8080/uploads/:path*'
          : 'http://backend:8080/uploads/:path*',
      },
    ];
  },
};

export default nextConfig;
