'use client'
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Poppins } from "next/font/google";
import { cn } from "@/lib/utils";
import { motion } from "framer-motion";
import { Terminal } from 'lucide-react'

const poppins = Poppins({
    weight: ["400", "600"],
    subsets: ["latin"],
    display: "swap",
});

export default function Navbar() {
    const pathname = usePathname();
    const navItems = [
        { href: "/", label: <Terminal /> },
        { href: "/work", label: "Work" },
        { href: "/blog", label: "Blog" },
    ];

    return (
        <motion.nav
            initial={{ y: -100, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.6, ease: "easeOut" }}
            className={cn(
                poppins.className,
                "flex justify-center absolute left-1/2 right-1/2 top-8 transform -translate-x-1/2"
            )}
        >
            <motion.ul
                initial={{ scale: 0.95 }}
                animate={{ scale: 1 }}
                transition={{ duration: 0.3, delay: 0.2 }}
                className="flex gap-x-8 px-0 py-1 text-xl border-2 rounded-2xl"
            >
                {navItems.map((item, index) => (
                    <motion.li
                        key={item.href}
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{
                            duration: 0.3,
                            delay: 0.3 + index * 0.1,
                            ease: "easeOut",
                        }}
                    >
                        <Link
                            href={item.href}
                            className={cn(
                                "px-4 py-1 mx-1 duration-500 ease-in-out transition-colors rounded-xl block relative",
                                pathname === item.href
                                    ? "bg-violet-400 text-white underline underline-offset-2"
                                    : "hover:bg-orange-100"
                            )}
                        >
                            {pathname === item.href && (
                                <motion.span
                                    layoutId="activeTab"
                                    className="bg-softPurple rounded-xl"
                                    transition={{ duration: 0.3 }}
                                />
                            )}
                            {item.label}
                        </Link>
                    </motion.li>
                ))}
            </motion.ul>
        </motion.nav>
    );
}