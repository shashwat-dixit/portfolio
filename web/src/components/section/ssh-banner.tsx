import { CONFIG } from "@/data/config";

export default function SshBanner() {
  return (
    <div className="rounded-xl border bg-card px-4 py-3">
      <p className="text-sm text-muted-foreground">Same site, in your terminal.</p>
      <p className="mt-1 font-mono text-sm">
        <span className="text-muted-foreground select-none">$ </span>
        <span>{CONFIG.site.sshCommand}</span>
      </p>
    </div>
  );
}
