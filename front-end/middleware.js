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

  const cookie = req.headers.get('cookie') || '';

  const authCheck = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/checkAuth`, {
    headers: {
      cookie,
    },
  });

  const isAuthenticated = authCheck.ok;
  const isProtectedRoute = protectedRoutes.some((route) =>
    path === route || path.startsWith(route + '/')
  );
  if (path === '/') {
    return NextResponse.redirect(
      new URL(isAuthenticated ? '/home' : '/login', req.url)
    );
  }

  // Redirect unauthenticated users away from protected pages
  if (!isAuthenticated && isProtectedRoute) {
    return NextResponse.redirect(new URL('/login', req.url));
  }

  // Redirect logged-in users away from login page
  if (isAuthenticated && (path === '/login' || path === '/register')) {
    return NextResponse.redirect(new URL('/home', req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/home',
    '/',
    '/profile',
    '/users/:path*',
    '/groups/:path*',
    '/groups',
    '/users',
    '/chat',
    '/login',
    '/register'
  ],
};
