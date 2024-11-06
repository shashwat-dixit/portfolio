import type { Metadata } from "next";
import "./globals.css";
import Navbar from "@/components/Navbar";

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
        <div className="fixed inset-0 -z-10">
          <div className="absolute inset-0 bg-gray-100 [background-size:32px_32px] [background-image:linear-gradient(to_right,rgb(15_25_42/0.2)_1px,transparent_1px),linear-gradient(to_bottom,rgb(15_25_42/0.1)_1px,transparent_1px)] dark:bg-gray-900 dark:opacity-20">
            <div className="absolute inset-0 bg-gradient-to-b from-white/60 to-white dark:from-gray-900/60 dark:to-gray-900" />
          </div>
        </div>
        <Navbar />
        {children}
      </body>
    </html>
  );
}
