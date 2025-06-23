// middleware.js
import { NextResponse } from 'next/server';

const protectedRoutes = [
  '/home',
  '/profile',
  '/group',
  '/groups',
  '/users',
  '/chat',
];

export async function middleware(req) {
  const url = req.nextUrl;
  const path = url.pathname;
  console.log(`[middleware] ${url.pathname}`);
  if (path.startsWith('/api') || path.startsWith('/uploads')) {
    return NextResponse.next();
  }

  const cookie = req.headers.get('cookie') || '';
  const baseUrl = process.env.NODE_ENV === 'development'
    ? 'http://127.0.0.1:8080'
    : 'http://backend:8080';
  const apiUrl = `${baseUrl}/api/checkAuth`;

  let isAuthenticated = false;

  try {
    const authCheck = await fetch(apiUrl, {
      headers: { cookie },
      cache: 'no-store',
    });
    isAuthenticated = authCheck.ok;
  } catch {
    isAuthenticated = false;
  }

  const isProtectedRoute = protectedRoutes.some(route =>
    path === route || path.startsWith(route + '/')
  );

  if (path === '/') {
    return NextResponse.redirect(new URL(isAuthenticated ? '/home' : '/login', req.url));
  }

  if (!isAuthenticated && isProtectedRoute) {
    return NextResponse.redirect(new URL('/login', req.url));
  }

  if (isAuthenticated && (path === '/login' || path === '/register')) {
    return NextResponse.redirect(new URL('/home', req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/',
    '/home',
    '/profile',
    '/group',
    '/groups/:path*',
    '/users/:path*',
    '/chat',
    '/login',
    '/register',
  ],
};