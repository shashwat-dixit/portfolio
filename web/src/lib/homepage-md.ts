import { CONFIG } from "@/data/config";
import { DATA } from "@/data/resume";

export function homepageMarkdown(): string {
  const origin = CONFIG.site.url.replace(/\/$/, "");
  const lines: string[] = [
    `# ${DATA.name}`,
    "",
    `> ${DATA.description}`,
    "",
    `${DATA.location}. [Website](${origin}) · [Email](mailto:${DATA.contact.email})`,
    "",
    "## About",
    "",
    DATA.summary,
    "",
    "## Work",
    "",
  ];

  for (const job of DATA.work) {
    const end = job.end ?? "Present";
    lines.push(`### [${job.company}](${job.href}) — ${job.title}`);
    lines.push("");
    lines.push(`${job.location}. ${job.start} – ${end}`);
    lines.push("");
    for (const bullet of job.bullets) {
      lines.push(`- ${bullet}`);
    }
    lines.push("");
  }

  lines.push("## Education", "");
  for (const school of DATA.education) {
    lines.push(`### [${school.school}](${school.href})`);
    lines.push("");
    lines.push(`${school.degree}. ${school.start} – ${school.end}`);
    lines.push("");
  }

  lines.push("## Skills", "");
  lines.push(DATA.skills.map((skill) => skill.name).join(", "));
  lines.push("", "## Projects", "");

  for (const project of DATA.projects) {
    const href = project.href || origin;
    lines.push(`### [${project.title}](${href})`);
    lines.push("");
    lines.push(project.description);
    if (project.technologies.length > 0) {
      lines.push("");
      lines.push(`Tech: ${project.technologies.join(", ")}`);
    }
    const links = project.links ?? [];
    if (links.length > 0) {
      lines.push("");
      for (const link of links) {
        lines.push(`- [${link.type}](${link.href})`);
      }
    }
    lines.push("");
  }

  lines.push("## Contact", "");
  lines.push(`- Email: ${DATA.contact.email}`);
  for (const social of Object.values(DATA.contact.social)) {
    if (social.name.toLowerCase().includes("email")) {
      continue;
    }
    lines.push(`- [${social.name}](${social.url})`);
  }
  lines.push("");
  lines.push(`Markdown index for agents: ${origin}/llms.txt`);
  lines.push("");
  return lines.join("\n");
}
