import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: 'Shashwat Dixit | Full-Stack Developer',
  description: 'Experienced Full-Stack Developer specializing in React, Next.js, Node.js, and Python. Skilled in building performant web applications and implementing innovative solutions.',
  keywords: ['Shashwat Dixit', 'Full-Stack Developer', 'React', 'Next.js', 'Node.js', 'Web Development', 'Bengaluru'],
  authors: [{ name: 'Shashwat Dixit' }],
  openGraph: {
    title: 'Shashwat Dixit | Full-Stack Developer',
    description: 'Experienced Full-Stack Developer specializing in React, Next.js, Node.js, and Python. Skilled in building performant web applications and implementing innovative solutions.',
    url: 'https://shashwatdixit.me',
    siteName: 'Shashwat Dixit Portfolio',
    images: [
      {
        url: '/images/profile-image.jpg',
        width: 1200,
        height: 630,
      },
    ],
    locale: 'en_US',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Shashwat Dixit | Full-Stack Developer',
    description: 'Experienced Full-Stack Developer specializing in React, Next.js, Node.js, and Python. Skilled in building performant web applications and implementing innovative solutions.',
    images: ['/images/profile-image.jpg'],
  },
  alternates: {
    canonical: 'https://shashwatdixit.me',
  }
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className="antialiased"
      >
        {children}
      </body>
    </html>
  );
}
