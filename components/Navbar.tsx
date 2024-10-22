import Link from "next/link"

const navLinks = [
    { href: "/", label: "Work" },
    { href: "/", label: "Contact" },
    { href: "/", label: "Blog" }
]

export default function Navbar() {
    return (
        <nav>
            <h1>
                <Link href="/">
                    Shashwat Dixit
                </Link>
            </h1>
            <ul>
                {navLinks.map((item) => (
                    <li key={item.href}>
                        <Link href={item.href}>
                            {item.label}
                        </Link>
                    </li>
                ))}
            </ul>
        </nav>
    )
}
